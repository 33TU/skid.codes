import { h } from "preact";
import { useContext } from "preact/hooks";
import { AuthContext } from "../../context";
import Code from "../../components/code";

const Paste = () => {
  const { authState } = useContext(AuthContext);

  return (
    <Code
      content={process.env.PREACT_APP_START_MSG}
      readonly={authState.authUser === undefined}
    />
  );
};

export default Paste;
