import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';
import { formatDate, formatFromBaseQuantity } from 'shared/formatters';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

import './PreApprovalRequest.css';

export function renderActionIcons(status, onEdit, onApproval, onDelete, shipmentAccessorialId) {
  // Only office users can approve requests.
  // If the request is approved/invoiced, they cannot be edited, only deleted.
  //TODO: hiding edit action until we have implementation
  return (
    <Fragment>
      {onApproval &&
        status === 'SUBMITTED' && (
          <span onClick={onApproval}>
            <FontAwesomeIcon className="icon actionable" icon={faCheck} />
          </span>
        )}
      {false && (
        <span onClick={onEdit}>
          <FontAwesomeIcon className="icon actionable" icon={faPencil} />
        </span>
      )}
      <span onClick={() => onDelete(shipmentAccessorialId)}>
        <FontAwesomeIcon className="icon actionable" icon={faTimes} />
      </span>
    </Fragment>
  );
}

const PreApprovalTable = ({ shipment_accessorials, isActionable, onEdit, onApproval, onDelete }) => (
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
            <tr key={row.id}>
              <td align="left">{row.accessorial.code}</td>
              <td align="left">{row.accessorial.item}</td>
              <td align="left"> {row.location[0].toUpperCase() + row.location.substring(1).toLowerCase()} </td>
              <td align="left">{formatFromBaseQuantity(row.quantity_1)}</td>
              <td align="left">{row.notes} </td>
              <td align="left">{formatDate(row.submitted_date)}</td>
              <td align="left">
                <span className="status">{status}</span>
                {row.status[0].toUpperCase() + row.status.substring(1).toLowerCase()}
              </td>
              <td>{isActionable && renderActionIcons(row.status, onEdit, onApproval, onDelete, row.id)}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  </div>
);

PreApprovalTable.propTypes = {
  shipment_accessorials: PropTypes.array,
  isActionable: PropTypes.bool,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onApproval: PropTypes.func,
};

export default PreApprovalTable;
