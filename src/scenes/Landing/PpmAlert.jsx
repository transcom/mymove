import React from 'react';
import Alert from 'shared/Alert';
import { ppmInfoPacket } from 'shared/constants';

const PpmAlert = (props) => {
  return (
    <Alert type="success" heading={props.heading}>
      Next, wait for approval. Once approved:
      <br />
      <ul>
        <li>
          Get certified <strong>weight tickets</strong>, both empty &amp; full
        </li>
        <li>
          Save <strong>expense receipts</strong>, including for storage
        </li>
        <li>
          Read the{' '}
          <strong>
            <a href={ppmInfoPacket} target="_blank" rel="noopener noreferrer" className="usa-link">
              PPM info sheet
            </a>
          </strong>{' '}
          for more info
        </li>
      </ul>
    </Alert>
  );
};

export default PpmAlert;
