import { h } from "preact";
import Code from "../../components/code";

interface Props {
  id: string;
  password?: string;
}

// Note: `user` comes from the URL, courtesy of our router
const Paste = ({ id, password }: Props) => {
  return <Code id={id} password={password} readonly={true} />;
};

export default Paste;
