import React from 'react';
import Marked from 'marked';
import './CertificationText.css';

function CertificationText({ certificationText }) {
  let certificationMarkup;
  if (certificationText) {
    certificationMarkup = Marked(certificationText);
  }

  return (
    <div className="certification_text_box">
      {certificationMarkup ? <div dangerouslySetInnerHTML={{ __html: certificationMarkup }} /> : 'Loading...'}
    </div>
  );
}

export default CertificationText;
