import React from 'react';
import PropTypes from 'prop-types';
const DocumentContent = props => {
  let content;
  if (props.contentType === 'application/pdf') {
    content = (
      <div className="pdf-placeholder">
        {props.filename && <span className="filename">{props.filename}</span>}
        This PDF can be <a href={props.url}>viewed here</a>.
      </div>
    );
  } else {
    content = (
      <img src={props.url} width="100%" height="100%" alt="document upload" />
    );
  }
  return <div className="page">{content}</div>;
};
DocumentContent.PropTypes = {
  contentType: PropTypes.string,
  filename: PropTypes.string,
  url: PropTypes.string.isRequired,
};

export default DocumentContent;
