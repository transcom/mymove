import React, { forwardRef, useImperativeHandle, useRef, useState } from 'react';

import DebounceButton from './DebounceButton';

export default {
  title: 'Components/DebounceButton',
  component: DebounceButton,
};

const FadeMessage = forwardRef((props, ref) => {
  const { text = 'Clicked!' } = props;
  const [visible, setVisible] = useState(false);
  const [style, setStyle] = useState({
    opacity: 0,
    transition: 'opacity 2s ease',
  });

  useImperativeHandle(ref, () => ({
    show() {
      setVisible(true);
      setStyle({ opacity: 1, transition: 'opacity 0s' });
      requestAnimationFrame(() => {
        setStyle({ opacity: 0, transition: 'opacity 1s ease' });
      });
      setTimeout(() => setVisible(false), 2000);
    },
  }));

  if (!visible) return null;

  return (
    <span
      style={{
        display: 'inline-block',
        ...style,
      }}
    >
      {text}
    </span>
  );
});

export const BasicDebounceButton = () => {
  const fadeRef = useRef(0);

  return (
    <div style={{ padding: 20, fontFamily: 'sans-serif' }}>
      <DebounceButton
        type="button"
        delay={3000}
        onClick={() => fadeRef.current?.show()}
        ariaLabel="button with debounce"
      >
        Three Second Debounce
      </DebounceButton>
      <FadeMessage ref={fadeRef} text="On Click Called" />
    </div>
  );
};
