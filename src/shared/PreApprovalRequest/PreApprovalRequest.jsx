import React, { Component, Fragment } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';
import { formatDate, formatFromBaseQuantity } from 'shared/formatters';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

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

export class PreApprovalRequest extends Component {
  state = { showDeleteForm: false };
  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded)
      if (this.state.closeOnSubmit) this.setState({ showForm: false });
      else this.props.clearForm();
  }
  onDelete = () => {
    this.setState({ showDeleteForm: true });
  };
  cancelDelete = () => {
    this.setState({ showDeleteForm: false });
  };
  render() {
    let row = this.props.shipmentLineItem;
    let status = '';
    if (isOfficeSite) {
      status = renderStatusIcon(row.status);
    }
    let deleteActiveClass = this.state.showDeleteForm ? 'delete-active' : '';
    return (
      <Fragment>
        <tr key={row.id} className={deleteActiveClass}>
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
          <td>
            {this.props.isActionable &&
              renderActionIcons(row.status, this.props.onEdit, this.props.onApproval, this.onDelete)}
          </td>
        </tr>
        {this.state.showDeleteForm && (
          <tr className="delete-confirm-row">
            <td colSpan="8" className="delete-confirm">
              <strong>Are you sure you want to delete?</strong>
              <button className="usa-button usa-button-secondary" onClick={this.cancelDelete}>
                No, do not delete
              </button>
              <button className="usa-button usa-button-secondary" onClick={this.props.onDelete}>
                Yes, delete
              </button>
            </td>
          </tr>
        )}
      </Fragment>
    );
  }
}
PreApprovalRequest.propTypes = {
  shipmentLineItem: PropTypes.object.isRequired,
  onEdit: PropTypes.func,
  onApproval: PropTypes.func,
  onDelete: PropTypes.func,
  isActionable: PropTypes.bool.isRequired,
};

export default PreApprovalRequest;
