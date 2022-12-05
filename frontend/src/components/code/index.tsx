import { h } from "preact";
import { useEffect, useState } from "preact/hooks";
import { fetchPaste } from "../../api/paste";

import CodeMirror from "@uiw/react-codemirror";
import { javascript } from "@codemirror/lang-javascript";
import { okaidia } from "@uiw/codemirror-theme-okaidia";

interface Props {
  id?: string;
  password?: string;
  content?: string;
  readonly?: boolean;
}

const Code = ({ id, password, content, readonly }: Props) => {
  const [code, setCode] = useState<string>(content ?? "Loading...");

  useEffect(() => {
    if (id) {
      (async () => {
        const res = await fetchPaste({
          id,
          password,
        });

        if (!res.content && res.password) {
          return setCode("Incorrect password.");
        }

        if (res.content) {
          return setCode(res.content);
        }
      })();
    }
  }, [id, content, password]);

  return (
    <div class="container">
      <CodeMirror
        value={code}
        readOnly={readonly}
        extensions={[javascript({ jsx: true })]}
        theme={okaidia}
        height="90vh"
      />
    </div>
  );
};

export default Code;
