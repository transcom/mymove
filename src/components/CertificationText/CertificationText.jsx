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
      <Box
        ref={scrollContainerRef}
        onScroll={handleScroll}
        sx={{
          mt: 2,
          maxHeight: 200,
          overflowY: 'auto',
          border: '1px solid #ccc',
          p: 2,
        }}
      >
        <div>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi vel urna in libero sollicitudin aliquam.
          Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Integer non turpis a
          ante scelerisque sodales. Vivamus malesuada libero non nunc luctus, et venenatis lorem semper. Praesent sit
          amet mi vel justo fermentum luctus. Suspendisse potenti. Lorem ipsum dolor sit amet, consectetur adipiscing
          elit. Morbi vel urna in libero sollicitudin aliquam. Vestibulum ante ipsum primis in faucibus orci luctus et
          ultrices posuere cubilia curae; Integer non turpis a ante scelerisque sodales. Vivamus malesuada libero non
          nunc luctus, et venenatis lorem semper. Praesent sit amet mi vel justo fermentum luctus. Suspendisse potenti.
        </div>
      </Box>
    </div>
  );
};

export default CertificationText;
