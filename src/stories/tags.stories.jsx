import React from 'react';
import { Tag } from '@trussworks/react-uswds';

import { faExclamationTriangle } from '@fortawesome/free-solid-svg-icons';

export default {
  title: 'Components|Tags',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/c5a114cc-845b-48ea-9b28-70f3b1110746?mode=design',
    },
  },
};

export const all = () => (
  <div id="tags" style={{ padding: '20px' }}>
    <hr />
    <h3>Tags</h3>
    <Tag>New</Tag>
    <Tag className="usa-tag--green">Authorized</Tag>
    <Tag className="usa-tag--red">Rejected</Tag>
    <Tag className="usa-tag--yellow">Pending</Tag>
    <Tag className="usa-tag--alert">
      <FontAwesomeIcon icon={faExclamationTriangle} />
    </Tag>
    <Tag className="usa-tag--teal">INTL</Tag>
    <Tag>3</Tag>
    <Tag className="usa-tag--cyan usa-tag--large">#ABC123K</Tag>
  </div>
);
