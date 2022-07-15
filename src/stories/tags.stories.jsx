import React from 'react';
import { Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export default {
  title: 'Components/Tags',
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
      <FontAwesomeIcon icon="exclamation" />
    </Tag>
    &nbsp;
    <Tag className="usa-tag--teal">INTL</Tag>
    <Tag>3</Tag>
    <Tag className="usa-tag--cyan usa-tag--large">#ABC123K</Tag>
  </div>
);
