import React, { useState, useRef, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { Box } from '@material-ui/core';

export const CertificationText = ({ certificationText, onScrollToBottom }) => {
  const [hasScrolledToBottom, setHasScrolledToBottom] = useState(false);
  const scrollContainerRef = useRef(null);
  const getTextMarkdown = (certification) => {
    let certificationMarkup;
    if (certification) {
      certificationMarkup = <ReactMarkdown>{certification}</ReactMarkdown>;
    }
    return certificationMarkup;
  };

  useEffect(() => {
    // Allow keyboard users to scroll via focus
    const el = scrollContainerRef.current;
    if (el) {
      el.setAttribute('tabindex', '0');
      el.setAttribute('role', 'region');
      el.setAttribute('aria-label', 'Agreement text');
    }
  }, []);

  const handleScroll = (e) => {
    const container = scrollContainerRef.current;

    if (container) {
      // Check if the user has scrolled to the bottom.
      const isAtBottom = Math.abs(e.target.scrollHeight - (e.target.scrollTop + e.target.clientHeight)) <= 1;
      if (isAtBottom) {
        if (!hasScrolledToBottom) {
          setHasScrolledToBottom(true);
          if (onScrollToBottom) {
            onScrollToBottom(true); // Notify the parent.
          }
        }
      }
    }
  };

  return (
    <div
      id="certificationTextScrollBox"
      data-testid="certificationTextScrollBox"
      className="certification_text_box"
      onScroll={handleScroll}
    >
      <Box data-testid="certificationTextBox" ref={scrollContainerRef} onScroll={handleScroll}>
        {getTextMarkdown(certificationText)}
      </Box>
    </div>
  );
};

export default CertificationText;
