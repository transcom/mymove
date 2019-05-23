import React, { Component, Fragment } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { capitalize } from 'lodash';
import moment from 'moment';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faBan from '@fortawesome/fontawesome-free-solid/faBan';
import faSignInAlt from '@fortawesome/fontawesome-free-solid/faSignInAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';

import './StorageInTransit.css';
import { formatDate4DigitYear } from 'shared/formatters';
import TspEditor from 'shared/StorageInTransit/TspEditor';
import OfficeEditor from 'shared/StorageInTransit/OfficeEditor';
import ApproveSitRequest from 'shared/StorageInTransit/ApproveSitRequest';
import DenySitRequest from 'shared/StorageInTransit/DenySitRequest';
import PlaceInSit from 'shared/StorageInTransit/PlaceInSit';
import { updateStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import { isOfficeSite, isTspSite } from 'shared/constants';
import SitStatusIcon from './SitStatusIcon';
import { sitDaysUsed } from 'shared/StorageInTransit/calculator';

export class StorageInTransit extends Component {
  constructor() {
    super();
    this.state = {
      showTspEditForm: false,
      showOfficeEditForm: false,
      showApproveForm: false,
      showDenyForm: false,
      showPlaceInSitForm: false,
      storageInTransit: {},
    };
  }

  componentDidMount() {
    this.authorizedStartDate();
  }

  authorizedStartDate = () => {
    const { storageInTransit } = this.props;
    return storageInTransit.authorized_start_date
      ? this.storageInTransitAuthorizedStartDate()
      : this.assignEstimatedStartDateToAuthorizedStartDate();
  };

  storageInTransitAuthorizedStartDate = () => {
    this.setState({
      storageInTransit: {
        ...this.props.storageInTransit,
      },
    });
  };

  assignEstimatedStartDateToAuthorizedStartDate = () => {
    this.setState({
      storageInTransit: {
        ...this.props.storageInTransit,
        authorized_start_date: this.props.storageInTransit.estimated_start_date,
      },
    });
  };

  openTspEditForm = () => {
    this.setState({ showTspEditForm: true });
  };

  closeTspEditForm = () => {
    this.setState({ showTspEditForm: false });
  };

  openOfficeEditForm = () => {
    this.setState({ showOfficeEditForm: true });
  };

  closeOfficeEditForm = () => {
    this.setState({ showOfficeEditForm: false });
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

  openPlaceInSitForm = () => {
    this.setState({ showPlaceInSitForm: true });
  };

  closePlaceInSitForm = () => {
    this.setState({ showPlaceInSitForm: false });
  };

  onSubmit = updatePayload => {
    this.props.updateStorageInTransit(
      this.props.storageInTransit.shipment_id,
      this.props.storageInTransit.id,
      updatePayload,
    );
  };

  renderSitActions() {
    if (isOfficeSite) {
      return this.renderOfficeSitActions();
    } else if (isTspSite) {
      return this.renderTspSitActions();
    } else {
      return null;
    }
  }

  renderOfficeSitActions() {
    const { storageInTransit } = this.props;

    switch (storageInTransit.status) {
      case 'REQUESTED':
        return (
          <Fragment>
            <span className="sit-action">
              <a onClick={this.openApproveForm}>
                <FontAwesomeIcon className="icon" icon={faCheck} />
                Approve
              </a>
            </span>
            <span className="sit-action">
              <a data-cy="deny-sit-link" onClick={this.openDenyForm}>
                <FontAwesomeIcon className="icon" icon={faBan} />
                Deny
              </a>
            </span>
          </Fragment>
        );
      case 'APPROVED':
      case 'DENIED':
        return (
          <span className="sit-action sit-edit" data-cy="sit-edit-link">
            <a onClick={this.openOfficeEditForm}>
              <FontAwesomeIcon className="icon" icon={faPencil} />
              Edit
            </a>
          </span>
        );
      // NOTE: No actions for IN_SIT, RELEASED, or DELIVERED statuses.
      default:
        return null;
    }
  }

  renderTspSitActions() {
    const { storageInTransit } = this.props;

    switch (storageInTransit.status) {
      case 'REQUESTED':
        return (
          <Fragment>
            <span className="sit-action sit-edit" data-cy="sit-edit-link">
              <a onClick={this.openTspEditForm}>
                <FontAwesomeIcon className="icon" icon={faPencil} />
                Edit
              </a>
            </span>
            <span data-cy="sit-delete-link" className="sit-action sit-delete">
              <a>
                <FontAwesomeIcon className="icon" icon={faTimes} />
                Delete
              </a>
            </span>
          </Fragment>
        );
      case 'APPROVED':
        return (
          <Fragment>
            <span className="sit-action place-in-sit">
              <a data-cy="place-in-sit-link" onClick={this.openPlaceInSitForm}>
                <FontAwesomeIcon className="icon" icon={faSignInAlt} />
                Place into SIT
              </a>
            </span>
            <span className="sit-action sit-delete">
              <a data-cy="sit-delete-link">
                <FontAwesomeIcon className="icon" icon={faTimes} />
                Delete
              </a>
            </span>
          </Fragment>
        );
      case 'DENIED':
        return (
          <span className="sit-action sit-delete">
            <a data-cy="sit-delete-link">
              <FontAwesomeIcon className="icon" icon={faTimes} />
              Delete
            </a>
          </span>
        );
      case 'IN_SIT':
        if (storageInTransit.location === 'ORIGIN') {
          // TODO: Release from SIT
          return (
            <span data-cy="sit-edit-link" className="sit-action sit-edit">
              <a onClick={this.openTspEditForm}>
                <FontAwesomeIcon className="icon" icon={faPencil} />
                Edit
              </a>
            </span>
          );
        } else if (storageInTransit.location === 'DESTINATION') {
          return (
            <Fragment>
              <span data-cy="sit-edit-link" className="sit-action sit-edit">
                <a onClick={this.openTspEditForm}>
                  <FontAwesomeIcon className="icon" icon={faPencil} />
                  Edit
                </a>
              </span>
              <span data-cy="sit-delete-link" className="sit-action sit-delete">
                <a>
                  <FontAwesomeIcon className="icon" icon={faTimes} />
                  Delete
                </a>
              </span>
            </Fragment>
          );
        } else {
          return null;
        }
      case 'RELEASED':
      case 'DELIVERED':
        return (
          <span data-cy="sit-edit-link" className="sit-action sit-edit">
            <a onClick={this.openTspEditForm}>
              <FontAwesomeIcon className="icon" icon={faPencil} />
              Edit
            </a>
          </span>
        );
      default:
        return null;
    }
  }

  render() {
    const { storageInTransit, daysRemaining } = this.props;
    const { showTspEditForm, showOfficeEditForm, showApproveForm, showDenyForm, showPlaceInSitForm } = this.state;
    const isDenied = storageInTransit.status === 'DENIED';
    const isRequested = storageInTransit.status === 'REQUESTED';
    const isApproved = storageInTransit.status === 'APPROVED';
    const isInSit = storageInTransit.status === 'IN_SIT';

    const daysUsed = sitDaysUsed(storageInTransit);

    return (
      <div data-cy="storage-in-transit" className="storage-in-transit">
        <div className="column-head">
          {capitalize(storageInTransit.location)} SIT
          <span className="unbold">
            {' '}
            <span className="sit-status-text" data-cy="sit-status-text">
              Status:
            </span>{' '}
            {isRequested && <SitStatusIcon isTspSite={isTspSite} />}
          </span>
          {isApproved ? (
            <span data-cy="storage-in-transit-status">
              <FontAwesomeIcon className="icon approval-ready" icon={faCheck} />
              Approved
            </span>
          ) : isDenied ? (
            <span data-cy="storage-in-transit-status-denied">
              <FontAwesomeIcon className="icon approval-problem" icon={faBan} />
              Denied
            </span>
          ) : isInSit ? (
            <span>In SIT{daysRemaining < 0 ? ' - SIT Expired' : ''}</span>
          ) : (
            <span>SIT {capitalize(storageInTransit.status)}</span>
          )}
          {showApproveForm ? (
            <ApproveSitRequest onClose={this.closeApproveForm} storageInTransit={this.state.storageInTransit} />
          ) : showDenyForm ? (
            <DenySitRequest onClose={this.closeDenyForm} storageInTransit={storageInTransit} />
          ) : showPlaceInSitForm ? (
            <PlaceInSit sit={storageInTransit} onClose={this.closePlaceInSitForm} />
          ) : showTspEditForm ? (
            <TspEditor
              updateStorageInTransit={this.onSubmit}
              onClose={this.closeTspEditForm}
              storageInTransit={storageInTransit}
            />
          ) : showOfficeEditForm ? (
            <OfficeEditor
              updateStorageInTransit={this.onSubmit}
              onClose={this.closeOfficeEditForm}
              storageInTransit={storageInTransit}
            />
          ) : (
            <span className="sit-actions">{this.renderSitActions()}</span>
          )}
        </div>
        {!showTspEditForm && (
          <div className="usa-width-one-whole">
            <div className="usa-width-one-half">
              <div className="sit-dates">
                <div className="column-subhead nested__same-font">Dates</div>
                <div className="panel-field nested__same-font">
                  <span className="field-title unbold">Est. start date</span>
                  <span className="field-value">{formatDate4DigitYear(storageInTransit.estimated_start_date)}</span>
                </div>
                {storageInTransit.actual_start_date && (
                  <div>
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">Actual start date</span>
                      <span className="field-value">{formatDate4DigitYear(storageInTransit.actual_start_date)}</span>
                    </div>
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">Days used</span>
                      <span data-cy="sit-days-used" className="field-value">
                        {daysUsed} days
                      </span>
                    </div>
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">Expires</span>
                      <span data-cy="sit-expires" className="field-value">
                        {isInSit
                          ? formatDate4DigitYear(
                              moment(storageInTransit.actual_start_date).add(daysRemaining + daysUsed, 'days'),
                            )
                          : 'na'}
                      </span>
                    </div>
                  </div>
                )}
              </div>
              {storageInTransit.notes !== undefined && (
                <div className="sit-notes">
                  <div className="column-subhead nested__same-font">Note</div>
                  <div className="panel-field nested__same-font">
                    <span className="field-title unbold">{storageInTransit.notes}</span>
                  </div>
                </div>
              )}
            </div>
            <div className="usa-width-one-half">
              {!isRequested && (
                <div className="sit-authorization-wrapper">
                  <div className="column-subhead nested__same-font">Authorization</div>
                  <div className="panel-field nested__same-font">
                    <span className="field-title unbold">SIT approved</span>
                    <span className="field-value">{isDenied ? 'No' : 'Yes'}</span>
                  </div>
                  {!isDenied && (
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">Earliest start date</span>
                      <span data-cy="sit-authorized-start-date" className="field-value">
                        {formatDate4DigitYear(storageInTransit.authorized_start_date)}
                      </span>
                    </div>
                  )}

                  {storageInTransit.authorization_notes && (
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">Note</span>
                      <span data-cy="sit-authorization-notes" className="field-value">
                        {storageInTransit.authorization_notes}
                      </span>
                    </div>
                  )}
                  {storageInTransit.sit_number && (
                    <div className="panel-field nested__same-font">
                      <span className="field-title unbold">SIT Number</span>
                      <span className="field-value">{storageInTransit.sit_number}</span>
                    </div>
                  )}
                </div>
              )}

              <div className="column-subhead nested__same-font">Warehouse</div>
              <div className="panel-field nested__same-font">
                <span className="field-title unbold">Warehouse ID</span>
                <span className="field-value">{storageInTransit.warehouse_id}</span>
              </div>
              <div className="panel-field nested__same-font">
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
  daysRemaining: PropTypes.number,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateStorageInTransit }, dispatch);
}

export default connect(null, mapDispatchToProps)(StorageInTransit);
