import { Editor } from "../../components/editor";
import { EditorView } from "@codemirror/view";

export default () => {
  const doc = import.meta.env.VITE_START_MSG;

  const init = (view: EditorView) => {
  };

  return (
    <div class="container-fluid">
      <Editor init={init} config={{ doc }}></Editor>
    </div>
  );
};
