import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { capitalize, get, includes } from 'lodash';

import { NavTab, RoutedTabs } from 'react-router-tabs';
import { Link, NavLink, Redirect, Switch } from 'react-router-dom';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'shared/User/PrivateRoute';
import LocationsContainer from './Hhg/LocationsContainer';
import Alert from 'shared/Alert'; // eslint-disable-line
import DocumentList from 'shared/DocumentViewer/DocumentList';
import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';
import PaymentsPanel from './Ppm/PaymentsPanel';
import DatesAndLocationPanel from './Ppm/DatesAndLocationsPanel';
import PPMEstimatesPanel from './Ppm/PPMEstimatesPanel';
import StoragePanel from './Ppm/StoragePanel';
import ExpensesPanel from './Ppm/ExpensesPanel';
import NetWeightPanel from './Ppm/NetWeightPanel';
import Dates from 'shared/ShipmentDates';
import RoutingPanel from './Hhg/RoutingPanel';
import ServiceAgentsContainer from './Hhg/ServiceAgentsContainer';
import Weights from 'shared/ShipmentWeights';
import PremoveSurvey from './PremoveSurvey';
import { withContext } from 'shared/AppContext';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';
import PreApprovalPanel from 'shared/PreApprovalRequest/PreApprovalPanel.jsx';
import StorageInTransitPanel from 'shared/StorageInTransit/StorageInTransitPanel.jsx';
import InvoicePanel from 'shared/Invoice/InvoicePanel.jsx';
import ComboButton from 'shared/ComboButton/index.jsx';
import ToolTip from 'shared/ToolTip';
import { DropDown, DropDownItem } from 'shared/ComboButton/dropdown';

import { getRequestStatus } from 'shared/Swagger/selectors';
import { resetRequests } from 'shared/Swagger/request';
import { getAllTariff400ngItems, selectTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';
import { getAllShipmentLineItems, selectSortedShipmentLineItems } from 'shared/Entities/modules/shipmentLineItems';
import { getAllInvoices } from 'shared/Entities/modules/invoices';
import { approvePPM, loadPPMs, selectPPMForMove, selectReimbursement } from 'shared/Entities/modules/ppms';
import { loadBackupContacts, loadServiceMember, selectServiceMember } from 'shared/Entities/modules/serviceMembers';
import { loadOrders, loadOrdersLabel, selectOrders } from 'shared/Entities/modules/orders';
import {
  approveShipment,
  completeShipment,
  getPublicShipment,
  selectShipment,
  selectShipmentStatus,
  updatePublicShipment,
} from 'shared/Entities/modules/shipments';
import { getTspForShipment } from 'shared/Entities/modules/transportationServiceProviders';
import { getServiceAgentsForShipment, selectServiceAgentsForShipment } from 'shared/Entities/modules/serviceAgents';

import { showBanner, removeBanner } from './ducks';
import {
  loadMove,
  loadMoveLabel,
  selectMove,
  selectMoveStatus,
  approveBasics,
  cancelMove,
} from 'shared/Entities/modules/moves';
import { formatDate } from 'shared/formatters';
import { getMoveDocumentsForMove, selectAllDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

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
      <OrdersPanel title="Orders" moveId={props.moveId} />
      <CustomerInfoPanel title="Customer Info" serviceMember={props.serviceMember} />
      <BackupInfoPanel title="Backup Info" serviceMember={props.serviceMember} />
      <AccountingPanel title="Accounting" serviceMember={props.serviceMember} moveId={props.moveId} />
    </div>
  );
};

const PPMTabContent = props => {
  return (
    <div className="office-tab">
      <PaymentsPanel title="Payments" moveId={props.moveId} />
      <ExpensesPanel title="Expenses" moveId={props.moveId} />
      <StoragePanel title="Storage" moveId={props.moveId} />
      <DatesAndLocationPanel title="Dates & Locations" moveId={props.moveId} />
      <NetWeightPanel title="Weights" moveId={props.moveId} />
      <PPMEstimatesPanel title="Estimates" moveId={props.moveId} />
    </div>
  );
};

const HHGTabContent = props => {
  let shipmentStatus = '';
  let shipmentId = '';
  const {
    allowHhgInvoicePayment,
    canApprovePaymentInvoice,
    moveId,
    serviceAgents,
    shipment,
    updatePublicShipment,
    showSitPanel,
  } = props;
  if (shipment) {
    shipmentStatus = shipment.status;
    shipmentId = shipment.id;
  }
  return (
    <div className="office-tab">
      <RoutingPanel title="Routing" moveId={moveId} />
      <Dates title="Dates" shipment={shipment} update={updatePublicShipment} />
      <LocationsContainer update={updatePublicShipment} shipmentId={shipment.id} />
      <Weights title="Weights & Items" shipment={shipment} update={updatePublicShipment} />
      <PremoveSurvey title="Premove Survey" shipment={shipment} update={updatePublicShipment} />
      <ServiceAgentsContainer
        title="TSP & Servicing Agents"
        shipment={shipment}
        serviceAgents={serviceAgents}
        transportationServiceProviderId={shipment.transportation_service_provider_id}
      />
      <PreApprovalPanel shipmentId={shipment.id} />
      {showSitPanel && <StorageInTransitPanel shipmentId={shipmentId} moveId={moveId} />}
      <InvoicePanel
        shipmentId={shipment.id}
        shipmentStatus={shipmentStatus}
        canApprove={canApprovePaymentInvoice}
        allowPayments={allowHhgInvoicePayment}
      />
    </div>
  );
};

class MoveInfo extends Component {
  state = {
    redirectToHome: false,
  };

  componentDidMount() {
    const { moveId } = this.props;
    this.props.loadMove(moveId);
    this.props.getMoveDocumentsForMove(moveId);
    this.props.getAllTariff400ngItems(true);
    this.props.loadPPMs(moveId);
  }

  componentDidUpdate(prevProps) {
    const {
      loadBackupContacts,
      loadOrders,
      loadMoveIsSuccess,
      loadServiceMember,
      ordersId,
      serviceMemberId,
      shipmentId,
    } = this.props;
    if (loadMoveIsSuccess !== prevProps.loadMoveIsSuccess && loadMoveIsSuccess) {
      loadOrders(ordersId);
      loadServiceMember(serviceMemberId);
      loadBackupContacts(serviceMemberId);
      if (shipmentId) {
        this.getAllShipmentInfo(shipmentId);
      }
    }
  }

  componentWillUnmount() {
    this.props.resetRequests();
  }

  getAllShipmentInfo = shipmentId => {
    this.props.getTspForShipment(shipmentId);
    this.props.getPublicShipment(shipmentId);
    this.props.getAllShipmentLineItems(shipmentId);
    this.props.getAllInvoices(shipmentId);
    this.props.getServiceAgentsForShipment(shipmentId);
  };

  approveBasics = () => {
    this.props.approveBasics(this.props.moveId);
  };

  approvePPM = () => {
    this.props.approvePPM(this.props.ppm.id);
  };

  approveShipment = () => {
    this.props.approveShipment(this.props.shipmentId);
  };

  completeShipment = () => {
    this.props.completeShipment(this.props.shipmentId);
  };

  cancelMoveAndRedirect = cancelReason => {
    const messageLines = [
      `Move #${this.props.move.locator} for ${this.props.serviceMember.last_name}, ${
        this.props.serviceMember.first_name
      } has been canceled`,
      'An email confirmation has been sent to the customer.',
    ];
    this.props.cancelMove(this.props.moveId, cancelReason).then(() => {
      this.props.showBanner({ messageLines });
      setTimeout(() => this.props.removeBanner(), 10000);
      this.setState({ redirectToHome: true });
    });
  };

  renderPPMTabStatus = () => {
    if (this.props.ppm.status === 'APPROVED') {
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
    const {
      move,
      moveDocuments,
      moveStatus,
      orders,
      ppm,
      shipment,
      shipmentStatus,
      serviceMember,
      upload,
    } = this.props;
    const isPPM = move.selected_move_type === 'PPM';
    const isHHG = move.selected_move_type === 'HHG';
    const isHHGPPM = move.selected_move_type === 'HHG_PPM';
    const pathnames = this.props.location.pathname.split('/');
    const currentTab = pathnames[pathnames.length - 1];
    const showDocumentViewer = this.props.context.flags.documentViewer;
    const moveInfoComboButton = this.props.context.flags.moveInfoComboButton;
    const allowHhgInvoicePayment = this.props.context.flags.allowHhgInvoicePayment;
    let check = <FontAwesomeIcon className="icon" icon={faCheck} />;
    const ordersComplete = Boolean(
      orders.orders_number && orders.orders_type_detail && orders.department_indicator && orders.tac,
    );
    const ppmApproved = includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    const hhgApproved = includes(['APPROVED', 'IN_TRANSIT', 'DELIVERED', 'COMPLETED'], shipmentStatus);
    const hhgAccepted = shipmentStatus === 'ACCEPTED';
    const hhgDelivered = shipmentStatus === 'DELIVERED';
    const hhgCompleted = shipmentStatus === 'COMPLETED';
    const moveApproved = moveStatus === 'APPROVED';
    const hhgCantBeCanceled = includes(['IN_TRANSIT', 'DELIVERED', 'COMPLETED'], shipmentStatus);

    const moveDate = isPPM ? ppm.original_move_date : shipment && shipment.requested_pickup_date;
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
                  <FontAwesomeIcon className="icon icon-grey" icon={faPhone} flip="horizontal" />
                )}
                {serviceMember.text_message_is_preferred && (
                  <FontAwesomeIcon className="icon icon-grey" icon={faComments} />
                )}
                {serviceMember.email_is_preferred && <FontAwesomeIcon className="icon icon-grey" icon={faEmail} />}
                &nbsp;
              </li>
              <li>Locator# {move.locator}&nbsp;</li>
              <li>Move date {formatDate(moveDate)}&nbsp;</li>
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
                  {capitalize(this.props.moveStatus)}
                </span>
              </NavTab>
              {(isPPM || isHHGPPM) && (
                <NavTab to="/ppm">
                  <span className="title">PPM</span>
                  {this.renderPPMTabStatus()}
                </NavTab>
              )}
              {(isHHG || isHHGPPM) && (
                <NavTab to="/hhg">
                  <span className="title">HHG</span>
                  <span className="status">
                    <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />
                    {capitalize(shipmentStatus)}
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
                <PrivateRoute path={`${this.props.match.path}/basics`}>
                  <BasicsTabContent moveId={this.props.moveId} serviceMember={this.props.serviceMember} />
                </PrivateRoute>
                <PrivateRoute path={`${this.props.match.path}/ppm`}>
                  <PPMTabContent moveId={this.props.moveId} />
                </PrivateRoute>
                <PrivateRoute path={`${this.props.match.path}/hhg`}>
                  {this.props.shipment && (
                    <HHGTabContent
                      allowHhgInvoicePayment={allowHhgInvoicePayment}
                      canApprovePaymentInvoice={hhgDelivered}
                      moveId={this.props.moveId}
                      serviceAgents={this.props.serviceAgents}
                      shipment={this.props.shipment}
                      shipmentStatus={this.props.shipmentStatus}
                      updatePublicShipment={this.props.updatePublicShipment}
                      showSitPanel={this.props.context.flags.sitPanel}
                    />
                  )}
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
              <div>
                <ToolTip
                  disabled={ordersComplete}
                  textStyle="tooltiptext-large"
                  toolTipText="Some information about the move is missing or contains errors. Please fix these problems before approving."
                >
                  {moveInfoComboButton && (
                    <ComboButton buttonText="Approve" disabled={!ordersComplete}>
                      <DropDown>
                        <DropDownItem
                          value="Approve Basics"
                          disabled={moveApproved || !ordersComplete}
                          onClick={this.approveBasics}
                        />
                        {(isPPM || isHHGPPM) && (
                          <DropDownItem
                            disabled={ppmApproved || !moveApproved || !ordersComplete}
                            onClick={this.approvePPM}
                            value="Approve PPM"
                          />
                        )}
                        {(isHHG || isHHGPPM) && (
                          <DropDownItem
                            value="Approve HHG"
                            onClick={this.approveShipment}
                            disabled={!hhgAccepted || hhgApproved || hhgCompleted || !moveApproved || !ordersComplete}
                          />
                        )}
                      </DropDown>
                    </ComboButton>
                  )}
                </ToolTip>
              </div>
              <button
                className={`${moveApproved ? 'btn__approve--green' : ''}`}
                onClick={this.approveBasics}
                disabled={moveApproved || !ordersComplete}
              >
                Approve Basics
                {moveApproved && check}
              </button>

              {(isPPM || isHHGPPM) && (
                <button
                  className={`${ppmApproved ? 'btn__approve--green' : ''}`}
                  onClick={this.approvePPM}
                  disabled={ppmApproved || !moveApproved || !ordersComplete}
                >
                  Approve PPM
                  {ppmApproved && check}
                </button>
              )}
              {(isHHG || isHHGPPM) && (
                <button
                  className={`${hhgApproved ? 'btn__approve--green' : ''}`}
                  onClick={this.approveShipment}
                  disabled={
                    !hhgAccepted ||
                    hhgApproved ||
                    hhgCompleted ||
                    !moveApproved ||
                    !ordersComplete ||
                    currentTab !== 'hhg'
                  }
                >
                  Approve HHG
                  {hhgApproved && check}
                </button>
              )}
              {(isHHG || isHHGPPM) && (
                <button
                  className={`${hhgCompleted ? 'btn__approve--green' : ''}`}
                  onClick={this.completeShipment}
                  disabled={!hhgDelivered || hhgCompleted || !moveApproved || !ordersComplete || currentTab !== 'hhg'}
                >
                  Complete Shipments
                  {hhgCompleted && check}
                </button>
              )}
              <ConfirmWithReasonButton
                buttonTitle="Cancel Move"
                reasonPrompt="Why is the move being canceled?"
                warningPrompt="Are you sure you want to cancel the entire move?"
                onConfirm={this.cancelMoveAndRedirect}
                buttonDisabled={hhgCantBeCanceled}
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
                <DocumentList detailUrlPrefix={`/moves/${this.props.moveId}/documents`} moveDocuments={moveDocuments} />
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
  loadMove: PropTypes.func.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({
      documentViewer: PropTypes.bool,
      sitPanel: PropTypes.bool,
    }).isRequired,
  }).isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const moveId = ownProps.match.params.moveId;
  const move = selectMove(state, moveId);
  const shipmentId = get(move, 'shipments.0.id');
  const ppm = selectPPMForMove(state, moveId);
  const ordersId = move.orders_id;
  const orders = selectOrders(state, ordersId);
  const serviceMemberId = move.service_member_id;
  const serviceMember = selectServiceMember(state, serviceMemberId);
  const loadOrdersStatus = getRequestStatus(state, loadOrdersLabel);
  const loadMoveIsSuccess = getRequestStatus(state, loadMoveLabel).isSuccess;

  return {
    approveMoveHasError: get(state, 'office.moveHasApproveError'),
    errorMessage: get(state, 'office.error'),
    loadDependenciesHasError: loadOrdersStatus.error,
    loadDependenciesHasSuccess: loadOrdersStatus.isSuccess,
    loadMoveIsSuccess,
    moveDocuments: selectAllDocumentsForMove(state, moveId),
    ppm,
    move,
    moveId,
    moveStatus: selectMoveStatus(state, moveId),
    orders,
    ordersId,
    officeShipment: get(state, 'office.officeShipment', {}),
    ppmAdvance: selectReimbursement(state, ppm.advance),
    serviceAgents: selectServiceAgentsForShipment(state, shipmentId),
    serviceMember,
    serviceMemberId,
    shipment: selectShipment(state, shipmentId),
    shipmentId,
    shipmentLineItems: selectSortedShipmentLineItems(state),
    shipmentStatus: selectShipmentStatus(state, shipmentId),
    swaggerError: get(state, 'swagger.hasErrored'),
    tariff400ngItems: selectTariff400ngItems(state),
    upload: get(orders, 'uploaded_orders.uploads.0', {}),
  };
};

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      getPublicShipment,
      updatePublicShipment,
      getMoveDocumentsForMove,
      approveBasics,
      approvePPM,
      approveShipment,
      completeShipment,
      cancelMove,
      getAllTariff400ngItems,
      getAllShipmentLineItems,
      getAllInvoices,
      getTspForShipment,
      getServiceAgentsForShipment,
      showBanner,
      removeBanner,
      loadMove,
      loadPPMs,
      loadServiceMember,
      loadBackupContacts,
      loadOrders,
      resetRequests,
    },
    dispatch,
  );

export default withContext(connect(mapStateToProps, mapDispatchToProps)(MoveInfo));
