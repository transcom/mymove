import React from 'react';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';
import { formatDate } from 'shared/formatters';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

import './PreApprovalRequest.css';

export function renderActionIcons(status, onEdit, onApproval, onDelete) {
  // Only office users can approve requests.
  // If the request is approved/invoiced, they cannot be edited, only deleted.
  if (status === 'APPROVED' || status === 'INVOICED') {
    return (
      <span>
        <span onClick={onDelete}>
          <FontAwesomeIcon className="icon actionable" icon={faTimes} />
        </span>
      </span>
    );
  } else if (onApproval) {
    if (status === 'SUBMITTED') {
      return (
        <span>
          <span onClick={onApproval}>
            <FontAwesomeIcon className="icon actionable" icon={faCheck} />
          </span>
          <span onClick={onEdit}>
            <FontAwesomeIcon className="icon actionable" icon={faPencil} />
          </span>
          <span onClick={onDelete}>
            <FontAwesomeIcon className="icon actionable" icon={faTimes} />
          </span>
        </span>
      );
    }
  } else {
    return (
      <span>
        <span onClick={onEdit}>
          <FontAwesomeIcon className="icon actionable" icon={faPencil} />
        </span>
        <span onClick={onDelete}>
          <FontAwesomeIcon className="icon actionable" icon={faTimes} />
        </span>
      </span>
    );
  }
}

const PreApprovalRequest = ({
  shipment_accessorials,
  isActionable,
  onEdit,
  onApproval,
  onDelete,
}) => (
  <div className="accessorial-panel">
    <table cellSpacing={0}>
      <tbody>
        <tr>
          <th>Code</th>
          <th>Item</th>
          <th>Loc.</th>
          <th>Base Quantity</th>
          <th>Notes</th>
          <th>Submitted</th>
          <th>Status</th>
          <th>&nbsp;</th>
        </tr>
        {shipment_accessorials.map(row => {
          let status = '';
          if (isOfficeSite) {
            status = renderStatusIcon(row.status);
          }
          return (
            <tr key={row.code}>
              <td align="left">{row.code}</td>
              <td align="left">{row.item}</td>
              <td align="left"> {row.location} </td>
              <td align="left">{row.base_quantity} </td>
              <td align="left">{row.notes} </td>
              <td align="left">{formatDate(row.created_at)}</td>
              <td align="left">
                <span className="status">{status}</span>
                {row.status}
              </td>
              <td>
                {isActionable &&
                  renderActionIcons(row.status, onEdit, onApproval, onDelete)}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  </div>
);

PreApprovalRequest.propTypes = {
  shipment_accessorials: PropTypes.array,
  isActionable: PropTypes.bool,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onApproval: PropTypes.func,
};

export default PreApprovalRequest;
