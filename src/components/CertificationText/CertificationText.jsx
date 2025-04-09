import React, { useState, useRef } from 'react';
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

  const handleScroll = () => {
    const container = scrollContainerRef.current;

    if (container) {
      // Check if the user has scrolled to the bottom.
      if (container.scrollTop + container.clientHeight >= container.scrollHeight) {
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
    <div className="certification_text_box" onScroll={handleScroll}>
      {getTextMarkdown(certificationText)}
    </div>
  );
};

export default CertificationText;
