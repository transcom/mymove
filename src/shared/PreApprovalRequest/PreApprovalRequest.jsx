import { get } from 'lodash';
import React, { Component, Fragment } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import PropTypes from 'prop-types';
import { withContext } from 'shared/AppContext';
import { renderStatusIcon } from 'shared/utils';
import { isOfficeSite } from 'shared/constants.js';
import { formatDate } from 'shared/formatters';
import Editor from 'shared/PreApprovalRequest/Editor.jsx';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';
import { isNewAccessorial } from 'shared/preApprovals';
import { getDetailsComponent } from './DetailsHelper';

function formatStatus(lineItem) {
  let formattedStatus = lineItem.status;
  if (lineItem.invoice_id) {
    formattedStatus = 'Invoiced';
  }

  return formattedStatus[0].toUpperCase() + formattedStatus.substring(1).toLowerCase();
}

export function renderActionIcons(status, onEdit, onApproval, onDelete, shipmentLineItemId) {
  // Only office users can approve requests.
  // If the request is approved/invoiced, they cannot be edited, only deleted.
  //TODO: hiding edit action until we have implementation
  return (
    <Fragment>
      {onApproval &&
        status === 'SUBMITTED' && (
          <span data-test="approve-request" onClick={() => onApproval(shipmentLineItemId)}>
            <FontAwesomeIcon className="icon actionable" icon={faCheck} />
          </span>
        )}
      {onEdit &&
        status === 'SUBMITTED' && (
          <span data-test="edit-request" onClick={onEdit}>
            <FontAwesomeIcon className="icon actionable" icon={faPencil} />
          </span>
        )}
      {onDelete && (
        <span data-test="delete-request" onClick={onDelete}>
          <FontAwesomeIcon className="icon actionable" icon={faTimes} />
        </span>
      )}
    </Fragment>
  );
}

export class PreApprovalRequest extends Component {
  state = { showDeleteForm: false, showEditForm: false };
  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      if (!this.props.isActionable && this.state.showDeleteForm) {
        this.cancelDelete();
      } else if (!this.props.isActionable && this.state.showEditForm) {
        this.cancelEdit();
      }
    }
  }
  onEdit = () => {
    this.props.isActive(true);
    this.setState({ showEditForm: true });
  };
  saveEdit = (shipmentLineItemId, editPayload) => {
    this.props.onEdit(shipmentLineItemId, editPayload);
  };
  cancelEdit = () => {
    this.props.isActive(false);
    this.setState({ showEditForm: false });
  };
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
    const hasInvoice = Boolean(row.invoice_id);
    const isShowingForm = Boolean(this.state.showDeleteForm || this.state.showEditForm);
    const showButtons = this.props.isActionable && !isShowingForm && !hasInvoice;
    const DetailsComponent = getDetailsComponent(
      row.tariff400ng_item.code,
      get(this.props, 'context.flags.robustAccessorial', false),
      isNewAccessorial(row),
    );
    if (this.state.showEditForm) {
      return (
        <tr>
          <td colSpan="8" className="pre-approval-form">
            <Editor
              tariff400ngItems={this.props.tariff400ngItems}
              shipmentLineItem={row}
              saveEdit={this.saveEdit}
              cancelEdit={this.cancelEdit}
              onSaveComplete={this.cancelEdit}
            />
          </td>
        </tr>
      );
    } else {
      let status = '';
      if (isOfficeSite && !hasInvoice) {
        status = renderStatusIcon(row.status);
      }
      const deleteActiveClass = this.state.showDeleteForm ? 'delete-active' : '';

      return (
        <Fragment>
          <tr key={row.id} className={deleteActiveClass} data-cy={row.tariff400ng_item.code}>
            <td>{row.tariff400ng_item.code}</td>
            <td>{row.tariff400ng_item.item}</td>
            <td> {row.location[0]} </td>
            <DetailsComponent {...this.props} />
            <td>{formatDate(row.submitted_date)}</td>
            <td>
              <span className="status">{status}</span>
              {formatStatus(row)}
            </td>
            <td>
              {showButtons && renderActionIcons(row.status, this.onEdit, this.props.onApproval, this.onDelete, row.id)}
            </td>
          </tr>
          {this.state.showDeleteForm && (
            <tr className="delete-confirm-row">
              <td colSpan="7" className="delete-confirm">
                <strong>Are you sure you want to delete?</strong>
                <button
                  className="usa-button usa-button-secondary"
                  onClick={this.cancelDelete}
                  data-test="cancel-delete"
                >
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
}
PreApprovalRequest.propTypes = {
  shipmentLineItem: PropTypes.object.isRequired,
  onEdit: PropTypes.func,
  onApproval: PropTypes.func,
  onDelete: PropTypes.func,
  isActive: PropTypes.func.isRequired,
  isActionable: PropTypes.bool.isRequired,
  tariff400ngItems: PropTypes.array,
};

export default withContext(PreApprovalRequest);
export { PreApprovalRequest as BasicPreApprovalRequest };
