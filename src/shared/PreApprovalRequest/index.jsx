import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

import './index.css';

export function renderActionIcons(status, props) {
  // Only office users can approve requests.
  // If the request is approved/invoiced, they cannot be edited.
  if (props.onApproval) {
    if (status === 'APPROVED' || status === 'INVOICED') {
      return (
        <span>
          <span onClick={props.onEdit}>
            <FontAwesomeIcon className="icon actionable" icon={faTimes} />
          </span>
        </span>
      );
    } else {
      return (
        <span>
          <span onClick={props.onApproval}>
            <FontAwesomeIcon className="icon actionable" icon={faCheck} />
          </span>
          <span onClick={props.onEdit}>
            <FontAwesomeIcon className="icon actionable" icon={faPencil} />
          </span>
          <span onClick={props.onDelete}>
            <FontAwesomeIcon className="icon actionable" icon={faTimes} />
          </span>
        </span>
      );
    }
  } else {
    if (status === 'APPROVED' || status === 'INVOICED') {
      return (
        <span>
          <span onClick={props.onDelete}>
            <FontAwesomeIcon className="icon actionable" icon={faTimes} />
          </span>
        </span>
      );
    } else {
      return (
        <span>
          <span onClick={props.onEdit}>
            <FontAwesomeIcon className="icon actionable" icon={faPencil} />
          </span>
          <span onClick={props.onDelete}>
            <FontAwesomeIcon className="icon actionable" icon={faTimes} />
          </span>
        </span>
      );
    }
  }
}

class PreApprovalRequest extends Component {
  render() {
    const { accessorials, isActionable } = this.props;
    return (
      <div className="accessorial-panel">
        <table cellSpacing={0}>
          <tbody>
            <tr>
              <th>Code</th>
              <th>Item</th>
              <th>Loc.</th>
              <th>Base Quantity</th>
              <th>Notes</th>
              <th>Status</th>
              <th>&nbsp;</th>
            </tr>
            {accessorials.map(row => {
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
                  <td align="left">
                    <span className="status">{status}</span>
                    {row.status}
                  </td>
                  <td>
                    {isActionable && renderActionIcons(row.status, this.props)}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  }
}

PreApprovalRequest.propTypes = {
  accessorials: PropTypes.array,
  isActionable: PropTypes.bool,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onApproval: PropTypes.func,
};

export default PreApprovalRequest;
