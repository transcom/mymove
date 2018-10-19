import React, { Component, Fragment } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import PropTypes from 'prop-types';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';
import { formatDate, formatFromBaseQuantity } from 'shared/formatters';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

export function renderActionIcons(status, onEdit, onApproval, onDelete, shipmentAccessorialId) {
  // Only office users can approve requests.
  // If the request is approved/invoiced, they cannot be edited, only deleted.
  //TODO: hiding edit action until we have implementation
  return (
    <Fragment>
      {onApproval &&
        status === 'SUBMITTED' && (
          <span data-test="approve-request" onClick={() => onApproval(shipmentAccessorialId)}>
            <FontAwesomeIcon className="icon actionable" icon={faCheck} />
          </span>
        )}
      {false && (
        <span data-test="edit-request" onClick={onEdit}>
          <FontAwesomeIcon className="icon actionable" icon={faPencil} />
        </span>
      )}
      {onDelete && (
        <span data-test="delete-request" onClick={() => onDelete()}>
          <FontAwesomeIcon className="icon actionable" icon={faTimes} />
        </span>
      )}
    </Fragment>
  );
}

export class PreApprovalRequest extends Component {
  state = { showDeleteForm: false };
  componentDidUpdate(prevProps, prevState, snapshot) {
    if (!this.props.isActionable && this.state.showDeleteForm) {
      this.cancelDelete();
    }
  }
  onDelete = () => {
    this.props.isActive(true);
    this.setState({ showDeleteForm: true });
  };
  cancelDelete = () => {
    this.props.isActive(false);
    this.setState({ showDeleteForm: false });
  };
  approveDelete = () => {
    this.props.isActive(false);
    this.props.onDelete(this.props.shipmentLineItem.id);
    // We don't want the user clicking delete more than once
    this.setState({ showDeleteForm: false });
  };
  render() {
    const row = this.props.shipmentLineItem;
    let status = '';
    if (isOfficeSite) {
      status = renderStatusIcon(row.status);
    }
    const deleteActiveClass = this.state.showDeleteForm ? 'delete-active' : '';
    const showButtons = this.props.isActionable && !this.state.showDeleteForm;
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
            {showButtons &&
              renderActionIcons(row.status, this.props.onEdit, this.props.onApproval, this.onDelete, row.id)}
          </td>
        </tr>
        {this.state.showDeleteForm && (
          <tr className="delete-confirm-row">
            <td colSpan="8" className="delete-confirm">
              <strong>Are you sure you want to delete?</strong>
              <button className="usa-button usa-button-secondary" onClick={this.cancelDelete} data-test="cancel-delete">
                No, do not delete
              </button>
              <button
                className="usa-button usa-button-secondary"
                onClick={this.approveDelete}
                data-test="approve-delete"
              >
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
  isActive: PropTypes.func.isRequired,
  isActionable: PropTypes.bool.isRequired,
};

export default PreApprovalRequest;
