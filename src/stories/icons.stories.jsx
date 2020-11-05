import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFile } from '@fortawesome/free-solid-svg-icons/faFile';
import { faCheckCircle } from '@fortawesome/free-regular-svg-icons/faCheckCircle';

// Icons
export default {
  title: 'Global|Icons',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/eabef4a2-603e-4c3d-b249-0e580f1c8306?mode=design',
    },
  },
};

export const all = () => (
  <div style={{ padding: '20px', background: '#f0f0f0' }}>
    <h3>Icons</h3>
    <div id="icons" style={{ display: 'flex', flexWrap: 'wrap' }}>
      <div>
        <FontAwesomeIcon icon={faFile} />
        <code>documents | faFile</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faCheckCircle} />
        <code>accept | faCheckCircle (regular)</code>
      </div>
    </div>
  </div>
);
