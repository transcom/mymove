import React from 'react';
import ReactMarkdown from 'react-markdown';
import './CertificationText.css';

function CertificationText({ certificationText }) {
  let certificationMarkup;
  if (certificationText) {
    certificationMarkup = <ReactMarkdown>{certificationText}</ReactMarkdown>;
  }

  return (
    <div className="certification_text_box">
      {certificationMarkup ? <div>{certificationMarkup}</div> : 'Loading...'}
    </div>
  );
}

export default CertificationText;
