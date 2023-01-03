import { useEffect, useId, useState } from "preact/hooks";

import { basicSetup } from "codemirror";
import { EditorState, EditorStateConfig } from "@codemirror/state";
import { EditorView } from "@codemirror/view";
import { oneDark } from "@codemirror/theme-one-dark";

import "./style.css";

export const defaultConfig: EditorStateConfig = {
  extensions: [basicSetup, oneDark],
};

export interface EditorProps {
  config?: EditorStateConfig;
  init?: (view: EditorView) => void;
  class?: string;
}

export const Editor = (props: EditorProps | undefined) => {
  const editorId = useId();
  const [editor, setEditor] = useState<EditorView>();

  useEffect(() => {
    // Try get element
    const element = window.document.getElementById(editorId);
    if (!element) {
      return;
    }

    // Create state
    const state = EditorState.create({
      ...defaultConfig,
      ...props?.config,
    });

    // Create view
    const view = new EditorView({
      state: state,
      parent: element,
    });

    // Set editor
    setEditor(view);

    // Call init
    if (props?.init) {
      props.init(view);
    }

    console.log("ass");
    // Unmount
    return () => {
      editor?.destroy();
    };
  }, []);

  return (
    <>
      <div class={props?.class} id={editorId} />
    </>
  );
};
