-- Extensions
create extension if not exists citext;
create extension if not exists pgcrypto;

-- Enums
create type user_role as enum (
    'guest', -- Non confirmed account
    'member', -- Confirmed account
    'admin' -- Admin account
);

-- Tables
create table users (
    id serial,
    username citext not null,
    email citext not null,
    password char(60) not null,

    role user_role default 'guest',
    created timestamp default now(),

    primary key(id),
    unique(username),
    unique(email),

    check(length(username) > 0 and length(username) <= 32),
    check(length(email) <= 256)
);

create table language_types (
    name citext, -- Language name
    mode text, -- Language highlight mode.
    mime text not null, -- Mime types.
    ext text[] not null, -- Possible extensions

    primary key (name)
);

create table users_pastes (
    id uuid default gen_random_uuid(),
    user_id int not null,
    lang_name citext not null,

    title text,

    password boolean not null,
    private boolean not null,
    unlisted boolean not null,

    content_text text, -- Used when password is false
    content_bytes bytea, -- Used when password is true

    created timestamp default now(),

    primary key(id),
    foreign key(user_id) references users(id),
    foreign key(lang_name) references language_types(name),

    check (
        length(title) <= 30
    ),
    check(
        length(content_text) <= 1024 * 1024 and
        length(content_bytes) <= 1024 * 1024
    ),
    check(
        -- Both can not bet set same time.
        content_text is not null and not password or -- Raw
        content_bytes is not null and password -- Crypted
    )
);

create table users_sessions (
    id serial,
    user_id int,

    addr inet not null,
    country char(2) not null,

    created timestamp not null default now(), -- When created.
    updated timestamp not null default now(), -- When last updated.
    revoked boolean not null default false, -- When revoked JWT can not anymore be refreshed.

    primary key(id),
    foreign key(user_id) references users(id)
);

-- Index for non duplicate active session per IP addr
create unique index users_sessions_revoked_check_cidx on users_sessions(user_id, addr, revoked) where not revoked;

-- Indexes for speeding up text based searches
create index users_username_idx on users(username collate "C");
create index users_pastes_lang_name_idx on users_pastes(lang_name collate "C");
create index users_pastes_title_idx on users_pastes(title collate "C");
create index users_pastes_content_text_idx on users_pastes(content_text collate "C");
create index users_pastes_content_text_vector on users_pastes using gin (to_tsvector('english', content_text));

-- Indexes for speeding up time based searches
create index users_created_idx on users(created);
create index users_pastes_idx on users_pastes(created);
create index users_sessions_created_idx on users_sessions(created);
create index users_sessions_updated_idx on users_sessions(updated);

-- Indexes for speeding up user id based searches
create index users_pastes_user_id_idx on users_pastes(user_id);
create index users_sessions_user_id_idx on users_sessions(user_id);

-- Non error throwing pgp decryption. Sets null on failure.
create or replace function pgp_sym_decrypt_null_on_err(data bytea, psw text) returns text AS $$
begin
  return pgp_sym_decrypt(data, psw);
exception
  when external_routine_invocation_exception then
    raise DEBUG using
       MESSAGE = format('Decryption failed: SQLSTATE %s, Msg: %s',
                        SQLSTATE,SQLERRM),
       HINT = 'pgp_sym_encrypt(...) failed; check your key',
       ERRCODE = 'external_routine_invocation_exception';
    return null;
end;
$$ language plpgsql;

-- Add admin account
insert into 
    users(username, email, password, role)
values
    ('admin', 'admin@localhost', crypt('dev12345', gen_salt('bf')), 'admin');

-- Add languages supported by CodeMirror as of 2022 december
insert into language_types(name, mode, mime, ext) values
   ('APL', 'apl', 'text/apl', array['dyalog', 'apl']::text[]),
   ('PGP', 'asciiarmor', 'undefined', array['asc', 'pgp', 'sig']::text[]),
   ('ASN.1', 'asn.1', 'text/x-ttcn-asn', array['asn', 'asn1']::text[]),
   ('Asterisk', 'asterisk', 'text/x-asterisk', array[]::text[]),
   ('Brainfuck', 'brainfuck', 'text/x-brainfuck', array['b', 'bf']::text[]),
   ('C', 'clike', 'text/x-csrc', array['c', 'h', 'ino']::text[]),
   ('C++', 'clike', 'text/x-c++src', array['cpp', 'c++', 'cc', 'cxx', 'hpp', 'h++', 'hh', 'hxx']::text[]),
   ('Cobol', 'cobol', 'text/x-cobol', array['cob', 'cpy', 'cbl']::text[]),
   ('C#', 'clike', 'text/x-csharp', array['cs']::text[]),
   ('Clojure', 'clojure', 'text/x-clojure', array['clj', 'cljc', 'cljx']::text[]),
   ('ClojureScript', 'clojure', 'text/x-clojurescript', array['cljs']::text[]),
   ('Closure Stylesheets (GSS)', 'css', 'text/x-gss', array['gss']::text[]),
   ('CMake', 'cmake', 'text/x-cmake', array['cmake', 'cmake.in']::text[]),
   ('CoffeeScript', 'coffeescript', 'undefined', array['coffee']::text[]),
   ('Common Lisp', 'commonlisp', 'text/x-common-lisp', array['cl', 'lisp', 'el']::text[]),
   ('Cypher', 'cypher', 'application/x-cypher-query', array['cyp', 'cypher']::text[]),
   ('Cython', 'python', 'text/x-cython', array['pyx', 'pxd', 'pxi']::text[]),
   ('Crystal', 'crystal', 'text/x-crystal', array['cr']::text[]),
   ('CSS', 'css', 'text/css', array['css']::text[]),
   ('CQL', 'sql', 'text/x-cassandra', array['cql']::text[]),
   ('D', 'd', 'text/x-d', array['d']::text[]),
   ('Dart', 'dart', 'undefined', array['dart']::text[]),
   ('diff', 'diff', 'text/x-diff', array['diff', 'patch']::text[]),
   ('Django', 'django', 'text/x-django', array[]::text[]),
   ('Dockerfile', 'dockerfile', 'text/x-dockerfile', array[]::text[]),
   ('DTD', 'dtd', 'application/xml-dtd', array['dtd']::text[]),
   ('Dylan', 'dylan', 'text/x-dylan', array['dylan', 'dyl', 'intr']::text[]),
   ('EBNF', 'ebnf', 'text/x-ebnf', array[]::text[]),
   ('ECL', 'ecl', 'text/x-ecl', array['ecl']::text[]),
   ('edn', 'clojure', 'application/edn', array['edn']::text[]),
   ('Eiffel', 'eiffel', 'text/x-eiffel', array['e']::text[]),
   ('Elm', 'elm', 'text/x-elm', array['elm']::text[]),
   ('Embedded JavaScript', 'htmlembedded', 'application/x-ejs', array['ejs']::text[]),
   ('Embedded Ruby', 'htmlembedded', 'application/x-erb', array['erb']::text[]),
   ('Erlang', 'erlang', 'text/x-erlang', array['erl']::text[]),
   ('Esper', 'sql', 'text/x-esper', array[]::text[]),
   ('Factor', 'factor', 'text/x-factor', array['factor']::text[]),
   ('FCL', 'fcl', 'text/x-fcl', array[]::text[]),
   ('Forth', 'forth', 'text/x-forth', array['forth', 'fth', '4th']::text[]),
   ('Fortran', 'fortran', 'text/x-fortran', array['f', 'for', 'f77', 'f90', 'f95']::text[]),
   ('F#', 'mllike', 'text/x-fsharp', array['fs']::text[]),
   ('Gas', 'gas', 'text/x-gas', array['s']::text[]),
   ('Gherkin', 'gherkin', 'text/x-feature', array['feature']::text[]),
   ('GitHub Flavored Markdown', 'gfm', 'text/x-gfm', array[]::text[]),
   ('Go', 'go', 'text/x-go', array['go']::text[]),
   ('Groovy', 'groovy', 'text/x-groovy', array['groovy', 'gradle']::text[]),
   ('HAML', 'haml', 'text/x-haml', array['haml']::text[]),
   ('Haskell', 'haskell', 'text/x-haskell', array['hs']::text[]),
   ('Haskell (Literate)', 'haskell-literate', 'text/x-literate-haskell', array['lhs']::text[]),
   ('Haxe', 'haxe', 'text/x-haxe', array['hx']::text[]),
   ('HXML', 'haxe', 'text/x-hxml', array['hxml']::text[]),
   ('ASP.NET', 'htmlembedded', 'application/x-aspx', array['aspx']::text[]),
   ('HTML', 'htmlmixed', 'text/html', array['html', 'htm', 'handlebars', 'hbs']::text[]),
   ('HTTP', 'http', 'message/http', array[]::text[]),
   ('IDL', 'idl', 'text/x-idl', array['pro']::text[]),
   ('Pug', 'pug', 'text/x-pug', array['jade', 'pug']::text[]),
   ('Java', 'clike', 'text/x-java', array['java']::text[]),
   ('Java Server Pages', 'htmlembedded', 'application/x-jsp', array['jsp']::text[]),
   ('JavaScript', 'javascript', 'undefined', array['js']::text[]),
   ('JSON', 'javascript', 'undefined', array['json', 'map']::text[]),
   ('JSON-LD', 'javascript', 'application/ld+json', array['jsonld']::text[]),
   ('JSX', 'jsx', 'text/jsx', array['jsx']::text[]),
   ('Jinja2', 'jinja2', 'text/jinja2', array['j2', 'jinja', 'jinja2']::text[]),
   ('Julia', 'julia', 'text/x-julia', array['jl']::text[]),
   ('Kotlin', 'clike', 'text/x-kotlin', array['kt']::text[]),
   ('LESS', 'css', 'text/x-less', array['less']::text[]),
   ('LiveScript', 'livescript', 'text/x-livescript', array['ls']::text[]),
   ('Lua', 'lua', 'text/x-lua', array['lua']::text[]),
   ('Markdown', 'markdown', 'text/x-markdown', array['markdown', 'md', 'mkd']::text[]),
   ('mIRC', 'mirc', 'text/mirc', array[]::text[]),
   ('MariaDB SQL', 'sql', 'text/x-mariadb', array[]::text[]),
   ('Mathematica', 'mathematica', 'text/x-mathematica', array['m', 'nb', 'wl', 'wls']::text[]),
   ('Modelica', 'modelica', 'text/x-modelica', array['mo']::text[]),
   ('MUMPS', 'mumps', 'text/x-mumps', array['mps']::text[]),
   ('MS SQL', 'sql', 'text/x-mssql', array[]::text[]),
   ('mbox', 'mbox', 'application/mbox', array['mbox']::text[]),
   ('MySQL', 'sql', 'text/x-mysql', array[]::text[]),
   ('Nginx', 'nginx', 'text/x-nginx-conf', array[]::text[]),
   ('NSIS', 'nsis', 'text/x-nsis', array['nsh', 'nsi']::text[]),
   ('NTriples', 'ntriples', 'undefined', array['nt', 'nq']::text[]),
   ('Objective-C', 'clike', 'text/x-objectivec', array['m']::text[]),
   ('Objective-C++', 'clike', 'text/x-objectivec++', array['mm']::text[]),
   ('OCaml', 'mllike', 'text/x-ocaml', array['ml', 'mli', 'mll', 'mly']::text[]),
   ('Octave', 'octave', 'text/x-octave', array['m']::text[]),
   ('Oz', 'oz', 'text/x-oz', array['oz']::text[]),
   ('Pascal', 'pascal', 'text/x-pascal', array['p', 'pas']::text[]),
   ('PEG.js', 'pegjs', 'null', array['jsonld']::text[]),
   ('Perl', 'perl', 'text/x-perl', array['pl', 'pm']::text[]),
   ('PHP', 'php', 'undefined', array['php', 'php3', 'php4', 'php5', 'php7', 'phtml']::text[]),
   ('Pig', 'pig', 'text/x-pig', array['pig']::text[]),
   ('Plain Text', 'null', 'text/plain', array['txt', 'text', 'conf', 'def', 'list', 'log']::text[]),
   ('PLSQL', 'sql', 'text/x-plsql', array['pls']::text[]),
   ('PostgreSQL', 'sql', 'text/x-pgsql', array[]::text[]),
   ('PowerShell', 'powershell', 'application/x-powershell', array['ps1', 'psd1', 'psm1']::text[]),
   ('Properties files', 'properties', 'text/x-properties', array['properties', 'ini', 'in']::text[]),
   ('ProtoBuf', 'protobuf', 'text/x-protobuf', array['proto']::text[]),
   ('Python', 'python', 'text/x-python', array['BUILD', 'bzl', 'py', 'pyw']::text[]),
   ('Puppet', 'puppet', 'text/x-puppet', array['pp']::text[]),
   ('Q', 'q', 'text/x-q', array['q']::text[]),
   ('R', 'r', 'text/x-rsrc', array['r', 'R']::text[]),
   ('reStructuredText', 'rst', 'text/x-rst', array['rst']::text[]),
   ('RPM Changes', 'rpm', 'text/x-rpm-changes', array[]::text[]),
   ('RPM Spec', 'rpm', 'text/x-rpm-spec', array['spec']::text[]),
   ('Ruby', 'ruby', 'text/x-ruby', array['rb']::text[]),
   ('Rust', 'rust', 'text/x-rustsrc', array['rs']::text[]),
   ('SAS', 'sas', 'text/x-sas', array['sas']::text[]),
   ('Sass', 'sass', 'text/x-sass', array['sass']::text[]),
   ('Scala', 'clike', 'text/x-scala', array['scala']::text[]),
   ('Scheme', 'scheme', 'text/x-scheme', array['scm', 'ss']::text[]),
   ('SCSS', 'css', 'text/x-scss', array['scss']::text[]),
   ('Shell', 'shell', 'undefined', array['sh', 'ksh', 'bash']::text[]),
   ('Sieve', 'sieve', 'application/sieve', array['siv', 'sieve']::text[]),
   ('Slim', 'slim', 'undefined', array['slim']::text[]),
   ('Smalltalk', 'smalltalk', 'text/x-stsrc', array['st']::text[]),
   ('Smarty', 'smarty', 'text/x-smarty', array['tpl']::text[]),
   ('Solr', 'solr', 'text/x-solr', array[]::text[]),
   ('SML', 'mllike', 'text/x-sml', array['sml', 'sig', 'fun', 'smackspec']::text[]),
   ('Soy', 'soy', 'text/x-soy', array['soy']::text[]),
   ('SPARQL', 'sparql', 'application/sparql-query', array['rq', 'sparql']::text[]),
   ('Spreadsheet', 'spreadsheet', 'text/x-spreadsheet', array[]::text[]),
   ('SQL', 'sql', 'text/x-sql', array['sql']::text[]),
   ('SQLite', 'sql', 'text/x-sqlite', array[]::text[]),
   ('Squirrel', 'clike', 'text/x-squirrel', array['nut']::text[]),
   ('Stylus', 'stylus', 'text/x-styl', array['styl']::text[]),
   ('Swift', 'swift', 'text/x-swift', array['swift']::text[]),
   ('sTeX', 'stex', 'text/x-stex', array[]::text[]),
   ('LaTeX', 'stex', 'text/x-latex', array['text', 'ltx', 'tex']::text[]),
   ('SystemVerilog', 'verilog', 'text/x-systemverilog', array['v', 'sv', 'svh']::text[]),
   ('Tcl', 'tcl', 'text/x-tcl', array['tcl']::text[]),
   ('Textile', 'textile', 'text/x-textile', array['textile']::text[]),
   ('TiddlyWiki', 'tiddlywiki', 'text/x-tiddlywiki', array[]::text[]),
   ('Tiki wiki', 'tiki', 'text/tiki', array[]::text[]),
   ('TOML', 'toml', 'text/x-toml', array['toml']::text[]),
   ('Tornado', 'tornado', 'text/x-tornado', array[]::text[]),
   ('troff', 'troff', 'text/troff', array['1', '2', '3', '4', '5', '6', '7', '8', '9']::text[]),
   ('TTCN', 'ttcn', 'text/x-ttcn', array['ttcn', 'ttcn3', 'ttcnpp']::text[]),
   ('TTCN_CFG', 'ttcn-cfg', 'text/x-ttcn-cfg', array['cfg']::text[]),
   ('Turtle', 'turtle', 'text/turtle', array['ttl']::text[]),
   ('TypeScript', 'javascript', 'application/typescript', array['ts']::text[]),
   ('TypeScript-JSX', 'jsx', 'text/typescript-jsx', array['tsx']::text[]),
   ('Twig', 'twig', 'text/x-twig', array[]::text[]),
   ('Web IDL', 'webidl', 'text/x-webidl', array['webidl']::text[]),
   ('VB.NET', 'vb', 'text/x-vb', array['vb']::text[]),
   ('VBScript', 'vbscript', 'text/vbscript', array['vbs']::text[]),
   ('Velocity', 'velocity', 'text/velocity', array['vtl']::text[]),
   ('Verilog', 'verilog', 'text/x-verilog', array['v']::text[]),
   ('VHDL', 'vhdl', 'text/x-vhdl', array['vhd', 'vhdl']::text[]),
   ('Vue.js Component', 'vue', 'undefined', array['vue']::text[]),
   ('XML', 'xml', 'undefined', array['xml', 'xsl', 'xsd', 'svg']::text[]),
   ('XQuery', 'xquery', 'application/xquery', array['xy', 'xquery']::text[]),
   ('Yacas', 'yacas', 'text/x-yacas', array['ys']::text[]),
   ('YAML', 'yaml', 'undefined', array['yaml', 'yml']::text[]),
   ('Z80', 'z80', 'text/x-z80', array['z80']::text[]),
   ('mscgen', 'mscgen', 'text/x-mscgen', array['mscgen', 'mscin', 'msc']::text[]),
   ('xu', 'mscgen', 'text/x-xu', array['xu']::text[]),
   ('msgenny', 'mscgen', 'text/x-msgenny', array['msgenny']::text[]),
   ('WebAssembly', 'wast', 'text/webassembly', array['wat', 'wast']::text[]);