import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faBan from '@fortawesome/fontawesome-free-solid/faBan';

import './StorageInTransit.css';
import { formatDate4DigitYear } from 'shared/formatters';
import Editor from 'shared/StorageInTransit/Editor';
import ApproveSitRequest from 'shared/StorageInTransit/ApproveSitRequest';
import DenySitRequest from 'shared/StorageInTransit/DenySitRequest';
import { updateStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import { isOfficeSite, isTspSite } from 'shared/constants.js';
import SitStatusIcon from './SitStatusIcon';

export class StorageInTransit extends Component {
  constructor() {
    super();
    this.state = {
      showEditForm: false,
      showApproveForm: false,
      showDenyForm: false,
      storageInTransit: {},
    };
  }

  componentDidMount() {
    this.addAuthorizedStartDate();
  }

  addAuthorizedStartDate = () => {
    this.setState({
      storageInTransit: {
        ...this.props.storageInTransit,
        authorized_start_date: this.props.storageInTransit.estimated_start_date,
      },
    });
  };

  openEditForm = () => {
    this.setState({ showEditForm: true });
  };

  closeEditForm = () => {
    this.setState({ showEditForm: false });
  };

  openApproveForm = () => {
    this.setState({ showApproveForm: true });
  };

  closeApproveForm = () => {
    this.setState({ showApproveForm: false });
  };

  openDenyForm = () => {
    this.setState({ showDenyForm: true });
  };

  closeDenyForm = () => {
    this.setState({ showDenyForm: false });
  };

  onSubmit = updatePayload => {
    this.props.updateStorageInTransit(
      this.props.storageInTransit.shipment_id,
      this.props.storageInTransit.id,
      updatePayload,
    );
  };

  render() {
    const { storageInTransit } = this.props;
    const { showEditForm, showApproveForm, showDenyForm } = this.state;

    return (
      <div className="storage-in-transit">
        <div className="column-head">
          {storageInTransit.location.charAt(0) + storageInTransit.location.slice(1).toLowerCase()} SIT
          <span className="unbold">
            {' '}
            <span className="sit-status-text">Status:</span>{' '}
            {storageInTransit.status === 'REQUESTED' && <SitStatusIcon isTspSite={isTspSite} />}
          </span>
          <span>SIT {storageInTransit.status.charAt(0) + storageInTransit.status.slice(1).toLowerCase()} </span>
          {showApproveForm ? (
            <ApproveSitRequest onClose={this.closeApproveForm} storageInTransit={this.state.storageInTransit} />
          ) : (
            isOfficeSite &&
            !showEditForm &&
            !showDenyForm && (
              <span className="sit-actions">
                <a className="approve-sit-link" onClick={this.openApproveForm}>
                  <FontAwesomeIcon className="icon" icon={faCheck} />
                  Approve
                </a>
              </span>
            )
          )}
          {showDenyForm ? (
            <DenySitRequest onClose={this.closeDenyForm} />
          ) : (
            isOfficeSite &&
            !showEditForm &&
            !showApproveForm && (
              <span className="sit-actions">
                <a className="deny-sit-link" onClick={this.openDenyForm}>
                  <FontAwesomeIcon className="icon" icon={faBan} />
                  Deny
                </a>
              </span>
            )
          )}
          {showEditForm ? (
            <Editor
              updateStorageInTransit={this.onSubmit}
              onClose={this.closeEditForm}
              storageInTransit={storageInTransit}
            />
          ) : isOfficeSite ? (
            <span className="sit-actions">
              <span className="sit-edit actionable">
                {storageInTransit.status === 'APPROVED' &&
                  !showApproveForm &&
                  !showDenyForm && (
                    <a onClick={this.openEditForm}>
                      <FontAwesomeIcon className="icon" icon={faPencil} />
                      Edit
                    </a>
                  )}
              </span>
            </span>
          ) : (
            <span className="sit-actions">
              <span className="sit-edit actionable">
                <a onClick={this.openEditForm}>
                  <FontAwesomeIcon className="icon" icon={faPencil} />
                  Edit
                </a>
              </span>
            </span>
          )}
        </div>
        {!showEditForm && (
          <div className="usa-width-one-whole">
            <div className="usa-width-one-half">
              <div className="column-subhead">Dates</div>
              <div className="panel-field">
                <span className="field-title unbold">Est. start date</span>
                <span className="field-value">{formatDate4DigitYear(storageInTransit.estimated_start_date)}</span>
              </div>
              {storageInTransit.notes !== undefined && (
                <div className="sit-notes">
                  <div className="column-subhead">Note</div>
                  <div className="panel-field">
                    <span className="field-title unbold">{storageInTransit.notes}</span>
                  </div>
                </div>
              )}
            </div>
            <div className="usa-width-one-half">
              <div className="column-subhead">Warehouse</div>
              <div className="panel-field">
                <span className="field-title unbold">Warehouse ID</span>
                <span className="field-value">{storageInTransit.warehouse_id}</span>
              </div>
              <div className="panel-field">
                <span className="field-title unbold">Contact info</span>
                <span className="field-value">
                  {storageInTransit.warehouse_name}
                  <br />
                  {storageInTransit.warehouse_address.street_address_1}
                  <br />
                  {storageInTransit.warehouse_address.street_address_2 && (
                    <span>
                      {storageInTransit.warehouse_address.street_address_2}
                      <br />
                    </span>
                  )}
                  {storageInTransit.warehouse_address.street_address_3 && (
                    <span>
                      {storageInTransit.warehouse_address.street_address_3}
                      <br />
                    </span>
                  )}
                  {storageInTransit.warehouse_address.city}, {storageInTransit.warehouse_address.state}{' '}
                  {storageInTransit.warehouse_address.postal_code}
                  {storageInTransit.warehouse_phone && (
                    <span>
                      <br />
                      {storageInTransit.warehouse_phone}
                    </span>
                  )}
                  {storageInTransit.warehouse_email && (
                    <span>
                      <br />
                      {storageInTransit.warehouse_email}
                    </span>
                  )}
                </span>
              </div>
            </div>
          </div>
        )}
      </div>
    );
  }
}

StorageInTransit.propTypes = {
  storageInTransit: PropTypes.object.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateStorageInTransit }, dispatch);
}

export default connect(null, mapDispatchToProps)(StorageInTransit);
