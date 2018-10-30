import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize, has, isEmpty, includes } from 'lodash';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { NavLink, Switch, Redirect, Link } from 'react-router-dom';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'shared/User/PrivateRoute';
import Alert from 'shared/Alert'; // eslint-disable-line
import DocumentList from 'shared/DocumentViewer/DocumentList';
import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';
import PaymentsPanel from './Ppm/PaymentsPanel';
import PPMEstimatesPanel from './Ppm/PPMEstimatesPanel';
import StorageReimbursementCalculator from './Ppm/StorageReimbursementCalculator';
import IncentiveCalculator from './Ppm/IncentiveCalculator';
import ExpensesPanel from './Ppm/ExpensesPanel';
import DatesAndTrackingPanel from './Hhg/DatesAndTrackingPanel';
import LocationsPanel from './Hhg/LocationsPanel';
import RoutingPanel from './Hhg/RoutingPanel';
import WeightAndInventoryPanel from './Hhg/WeightAndInventoryPanel';
import ServiceAgents from './ServiceAgents';
import PremoveSurvey from 'shared/PremoveSurvey';
import { withContext } from 'shared/AppContext';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';
import PreApprovalPanel from 'shared/PreApprovalRequest/PreApprovalPanel.jsx';
import {
  getAllTariff400ngItems,
  selectTariff400ngItems,
  getTariff400ngItemsLabel,
} from 'shared/Entities/modules/tariff400ngItems';
import {
  getAllShipmentLineItems,
  selectShipmentLineItems,
  getShipmentLineItemsLabel,
} from 'shared/Entities/modules/shipmentLineItems';

import {
  loadMoveDependencies,
  approveBasics,
  approvePPM,
  approveHHG,
  completeHHG,
  cancelMove,
  patchShipment,
  sendHHGInvoice,
  resetMove,
} from './ducks';
import { formatDate } from 'shared/formatters';
import { selectAllDocumentsForMove, getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

const BasicsTabContent = props => {
  return (
    <div className="office-tab">
      <OrdersPanel title="Orders" />
      <CustomerInfoPanel title="Customer Info" moveId={props.match.params.moveId} />
      <BackupInfoPanel title="Backup Info" moveId={props.match.params.moveId} />
      <AccountingPanel title="Accounting" moveId={props.match.params.moveId} />
    </div>
  );
};

const PPMTabContent = props => {
  return (
    <div className="office-tab">
      <PaymentsPanel title="Payments" moveId={props.match.params.moveId} />
      <ExpensesPanel title="Expenses" />
      <IncentiveCalculator />
      <StorageReimbursementCalculator />
      <PPMEstimatesPanel title="Estimates" moveId={props.match.params.moveId} />
    </div>
  );
};

const HHGTabContent = props => {
  return (
    <div className="office-tab">
      <RoutingPanel title="Routing" moveId={props.moveId} />
      <DatesAndTrackingPanel title="Dates & Tracking" moveId={props.moveId} />
      <LocationsPanel title="Locations" moveId={props.moveId} />
      <WeightAndInventoryPanel title="Weight & Inventory" moveId={props.moveId} />
      {props.officeShipment && (
        <PremoveSurvey
          title="Premove Survey"
          shipment={props.officeShipment}
          update={props.patchShipment}
          error={props.surveyError}
        />
      )}
      {props.officeShipment.service_agents && (
        <ServiceAgents
          title="Service Agents"
          shipment={props.officeShipment}
          serviceAgents={props.officeShipment.service_agents}
        />
      )}
      {has(props, 'officeShipment.id') && <PreApprovalPanel shipmentId={props.officeShipment.id} />}
    </div>
  );
};

class MoveInfo extends Component {
  state = {
    redirectToHome: false,
  };

  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
    this.props.getAllTariff400ngItems(true, getTariff400ngItemsLabel);
  }

  componentDidUpdate(prevProps) {
    if (get(this.props, 'officeShipment.id') !== get(prevProps, 'officeShipment.id')) {
      this.props.getAllShipmentLineItems(getShipmentLineItemsLabel, this.props.officeShipment.id);
    }
  }

  componentWillUnmount() {
    this.props.resetMove();
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.match.params.moveId);
  };

  approvePPM = () => {
    this.props.approvePPM(this.props.officeMove.id, this.props.officePPM.id);
  };

  approveHHG = () => {
    this.props.approveHHG(this.props.officeShipment.id);
  };

  completeHHG = () => {
    this.props.completeHHG(this.props.officeShipment.id);
  };

  submitInvoice = () => {
    this.props.sendHHGInvoice(this.props.officeShipment.id);
  };

  cancelMove = cancelReason => {
    this.props.cancelMove(this.props.officeMove.id, cancelReason).then(() => {
      this.setState({ redirectToHome: true });
    });
  };

  renderPPMTabStatus = () => {
    if (this.props.officePPM.status === 'APPROVED') {
      if (this.props.ppmAdvance.status === 'APPROVED' || !this.props.ppmAdvance.status) {
        return (
          <span className="status">
            <FontAwesomeIcon className="icon approval-ready" icon={faCheck} />
            Move pending
          </span>
        );
      } else {
        return (
          <span className="status">
            <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />
            Payment Requested
          </span>
        );
      }
    } else {
      return (
        <span className="status">
          <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />
          In review
        </span>
      );
    }
  };

  render() {
    const { moveDocuments } = this.props;
    const move = this.props.officeMove;
    const serviceMember = this.props.officeServiceMember;
    const orders = this.props.officeOrders;
    const ppm = this.props.officePPM;
    const hhg = this.props.officeHHG;
    const isPPM = !isEmpty(this.props.officePPM);
    const isHHG = !isEmpty(this.props.officeHHG);
    const pathnames = this.props.location.pathname.split('/');
    const invoiceSuccess = this.props.hhgInvoiceHasSendSuccess;
    const currentTab = pathnames[pathnames.length - 1];

    const showDocumentViewer = this.props.context.flags.documentViewer;
    let upload = get(this.props, 'officeOrders.uploaded_orders.uploads.0'); // there can be only one
    let check = <FontAwesomeIcon className="icon" icon={faCheck} />;
    const ordersComplete = Boolean(
      orders.orders_number && orders.orders_type_detail && orders.department_indicator && orders.tac,
    );
    const ppmApproved = includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    const hhgApproved = includes(['APPROVED', 'IN_TRANSIT', 'DELIVERED', 'COMPLETED'], hhg.status);
    const hhgAccepted = hhg.status === 'ACCEPTED';
    const hhgDelivered = hhg.status === 'DELIVERED';
    const hhgCompleted = hhg.status === 'COMPLETED';
    const moveApproved = move.status === 'APPROVED';
    if (this.state.redirectToHome) {
      return <Redirect to="/" />;
    }

    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="usa-grid">
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              Something went wrong contacting the server.
            </Alert>
          </div>
        </div>
      );

    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds">
            <h1>
              Move Info: {serviceMember.last_name}, {serviceMember.first_name}
            </h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole">
            <ul className="move-info-header-meta">
              <li>ID# {serviceMember.edipi}&nbsp;</li>
              <li>
                {serviceMember.telephone}
                {serviceMember.phone_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faPhone} flip="horizontal" />
                )}
                {serviceMember.text_message_is_preferred && <FontAwesomeIcon className="icon" icon={faComments} />}
                {serviceMember.email_is_preferred && <FontAwesomeIcon className="icon" icon={faEmail} />}
                &nbsp;
              </li>
              <li>Locator# {move.locator}&nbsp;</li>
              <li>Move date {formatDate(ppm.planned_move_date)}&nbsp;</li>
            </ul>
          </div>
        </div>

        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-three-fourths">
            <RoutedTabs startPathWith={this.props.match.url}>
              <NavTab to="/basics">
                <span className="title">Basics</span>
                <span className="status">
                  <FontAwesomeIcon className="icon" icon={faPlayCircle} />
                  {capitalize(move.status)}
                </span>
              </NavTab>
              {isPPM && (
                <NavTab to="/ppm">
                  <span className="title">PPM</span>
                  {this.renderPPMTabStatus()}
                </NavTab>
              )}
              {isHHG && (
                <NavTab to="/hhg">
                  <span className="title">HHG</span>
                  <span className="status">
                    <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />
                    {capitalize(hhg.status)}
                  </span>
                </NavTab>
              )}
            </RoutedTabs>

            <div className="tab-content">
              <Switch>
                <PrivateRoute
                  exact
                  path={`${this.props.match.url}`}
                  render={() => <Redirect replace to={`${this.props.match.url}/basics`} />}
                />
                <PrivateRoute path={`${this.props.match.path}/basics`} component={BasicsTabContent} />
                <PrivateRoute path={`${this.props.match.path}/ppm`} component={PPMTabContent} />
                <PrivateRoute path={`${this.props.match.path}/hhg`}>
                  <HHGTabContent
                    officeHHG={JSON.stringify(this.props.officeHHG)}
                    officeShipment={this.props.officeShipment}
                    patchShipment={this.props.patchShipment}
                    moveId={this.props.match.params.moveId}
                    surveyError={this.props.shipmentPatchError && this.props.errorMessage}
                  />
                </PrivateRoute>
              </Switch>
            </div>
          </div>
          <div className="usa-width-one-fourth">
            <div>
              {this.props.approveMoveHasError && (
                <Alert type="warning" heading="Unable to approve">
                  Please fill out missing data
                </Alert>
              )}
              {this.props.hhgInvoiceHasSendSuccess && (
                <Alert type="success" heading="Success">
                  Invoice successfully sent
                </Alert>
              )}
              {this.props.hhgInvoiceHasFailure && (
                <Alert type="error" heading="">
                  Unable to send invoice. Please try again in a few minutes.
                </Alert>
              )}
              <button
                className={`${moveApproved ? 'btn__approve--green' : ''}`}
                onClick={this.approveBasics}
                disabled={moveApproved || !ordersComplete}
              >
                Approve Basics
                {moveApproved && check}
              </button>
              {isPPM ? (
                <button
                  className={`${ppmApproved ? 'btn__approve--green' : ''}`}
                  onClick={this.approvePPM}
                  disabled={ppmApproved || !moveApproved || !ordersComplete}
                >
                  Approve PPM
                  {ppmApproved && check}
                </button>
              ) : (
                <button
                  className={`${hhgApproved ? 'btn__approve--green' : ''}`}
                  onClick={this.approveHHG}
                  disabled={
                    !hhgAccepted ||
                    hhgApproved ||
                    hhgCompleted ||
                    !moveApproved ||
                    !ordersComplete ||
                    currentTab !== 'hhg'
                  }
                >
                  Approve Shipments
                  {hhgApproved && check}
                </button>
              )}
              {isHHG && (
                <button
                  className={`${hhgCompleted ? 'btn__approve--green' : ''}`}
                  onClick={this.completeHHG}
                  disabled={!hhgDelivered || hhgCompleted || !moveApproved || !ordersComplete || currentTab !== 'hhg'}
                >
                  Complete Shipments
                  {hhgCompleted && check}
                </button>
              )}
              <button
                className={`${invoiceSuccess ? 'btn__approve--green' : ''}`}
                onClick={this.submitInvoice}
                disabled={
                  !hhgCompleted ||
                  !hhgApproved ||
                  !moveApproved ||
                  !ordersComplete ||
                  invoiceSuccess ||
                  currentTab !== 'hhg'
                }
              >
                Submit HHG Invoice
                {invoiceSuccess && check}
              </button>

              <ConfirmWithReasonButton
                buttonTitle="Cancel Move"
                reasonPrompt="Why is the move being canceled?"
                warningPrompt="Are you sure you want to cancel the entire move?"
                onConfirm={this.cancelMove}
              />
              {/* Disabling until features implemented
              <button>Troubleshoot</button>
              */}
            </div>
            <div className="documents">
              <h2 className="extras usa-heading">
                Documents
                {!showDocumentViewer && <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />}
                {showDocumentViewer && (
                  <Link to={`/moves/${move.id}/documents`} target="_blank">
                    <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
                  </Link>
                )}
              </h2>
              {!upload ? (
                <p>No orders have been uploaded.</p>
              ) : (
                <div>
                  {moveApproved ? (
                    <div className="panel-field">
                      <FontAwesomeIcon style={{ color: 'green' }} className="icon" icon={faCheck} />
                      <Link to={`/moves/${move.id}/orders`} target="_blank">
                        Orders ({formatDate(upload.created_at)})
                      </Link>
                    </div>
                  ) : (
                    <div className="panel-field">
                      <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
                      <Link to={`/moves/${move.id}/orders`} target="_blank">
                        Orders ({formatDate(upload.created_at)})
                      </Link>
                    </div>
                  )}
                </div>
              )}
              {showDocumentViewer && (
                <DocumentList
                  detailUrlPrefix={`/moves/${this.props.match.params.moveId}/documents`}
                  moveDocuments={moveDocuments}
                />
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

MoveInfo.defaultProps = {
  move: {},
};

MoveInfo.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({ documentViewer: PropTypes.bool }).isRequired,
  }).isRequired,
};

const mapStateToProps = state => ({
  swaggerError: get(state, 'swagger.hasErrored'),
  officeMove: get(state, 'office.officeMove', {}),
  officeShipment: get(state, 'office.officeShipment', {}),
  officeOrders: get(state, 'office.officeOrders', {}),
  officeServiceMember: get(state, 'office.officeServiceMember', {}),
  officeBackupContacts: get(state, 'office.officeBackupContacts', []),
  officePPM: get(state, 'office.officePPMs.0', {}),
  officeHHG: get(state, 'office.officeMove.shipments.0', {}),
  ppmAdvance: get(state, 'office.officePPMs.0.advance', {}),
  moveDocuments: selectAllDocumentsForMove(state, get(state, 'office.officeMove.id', '')),
  tariff400ngItems: selectTariff400ngItems(state),
  shipmentLineItems: selectShipmentLineItems(state),
  loadDependenciesHasSuccess: get(state, 'office.loadDependenciesHasSuccess'),
  loadDependenciesHasError: get(state, 'office.loadDependenciesHasError'),
  shipmentPatchError: get(state, 'office.shipmentPatchError'),
  approveMoveHasError: get(state, 'office.moveHasApproveError'),
  hhgInvoiceHasSendSuccess: get(state, 'office.hhgInvoiceHasSendSuccess'),
  hhgInvoiceHasFailure: get(state, 'office.hhgInvoiceHasFailure'),
  errorMessage: get(state, 'office.error'),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadMoveDependencies,
      getMoveDocumentsForMove,
      approveBasics,
      approvePPM,
      approveHHG,
      completeHHG,
      cancelMove,
      patchShipment,
      sendHHGInvoice,
      getAllTariff400ngItems,
      getAllShipmentLineItems,
      resetMove,
    },
    dispatch,
  );

export default withContext(connect(mapStateToProps, mapDispatchToProps)(MoveInfo));
