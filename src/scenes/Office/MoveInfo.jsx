import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { capitalize, get, includes } from 'lodash';
import { NavTab, RoutedTabs } from 'react-router-tabs';
import { NavLink, Redirect, Switch } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import moment from 'moment';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'containers/PrivateRoute';
import Alert from 'shared/Alert';
import ToolTip from 'shared/ToolTip';
import ComboButton from 'shared/ComboButton';
import { DropDown, DropDownItem } from 'shared/ComboButton/dropdown';
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
import WeightsPanel from './Ppm/WeightsPanel';
import { withContext } from 'shared/AppContext';
import ConfirmWithReasonButton from 'shared/ConfirmWithReasonButton';

import { getRequestStatus } from 'shared/Swagger/selectors';
import { resetRequests } from 'shared/Swagger/request';
import { approvePPM, loadPPMs, selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import { loadBackupContacts, loadServiceMember, selectServiceMember } from 'shared/Entities/modules/serviceMembers';
import { loadOrders, loadOrdersLabel, selectOrders } from 'shared/Entities/modules/orders';
import { selectReimbursementById } from 'store/entities/selectors';
import { openLinkInNewWindow } from 'shared/utils';
import { defaultRelativeWindowSize } from 'shared/constants';

import { roleTypes } from 'constants/userRoles';

import { showBanner, removeBanner } from './ducks';
import {
  loadMove,
  loadMoveLabel,
  selectMove,
  selectMoveStatus,
  approveBasics,
  cancelMove,
} from 'shared/Entities/modules/moves';
import { formatDate } from 'utils/formatters';
import { getMoveDocumentsForMove, selectAllDocumentsForMove } from 'shared/Entities/modules/moveDocuments';

const BasicsTabContent = (props) => {
  return (
    <div className="office-tab">
      <OrdersPanel title="Orders" moveId={props.moveId} />
      <CustomerInfoPanel title="Customer Info" serviceMember={props.serviceMember} />
      <BackupInfoPanel title="Backup Info" serviceMember={props.serviceMember} />
      <AccountingPanel title="Accounting" serviceMember={props.serviceMember} moveId={props.moveId} />
    </div>
  );
};

const PPMTabContent = (props) => {
  return (
    <div className="office-tab">
      <PaymentsPanel title="Payments" moveId={props.moveId} />
      {props.ppmPaymentRequested && (
        <>
          <ExpensesPanel title="Expenses" moveId={props.moveId} moveDocuments={props.moveDocuments} />
          <StoragePanel title="Storage" moveId={props.moveId} moveDocuments={props.moveDocuments} />
          <DatesAndLocationPanel title="Dates & Locations" moveId={props.moveId} />
          <WeightsPanel title="Weights" moveId={props.moveId} ppmPaymentRequestedFlag={props.ppmPaymentRequestedFlag} />
        </>
      )}

      <PPMEstimatesPanel title="Estimates" moveId={props.moveId} />
    </div>
  );
};

const ReferrerQueueLink = (props) => {
  const pathname = props.history.location.state ? props.history.location.state.referrerPathname : '';
  switch (pathname) {
    case '/queues/ppm':
      return (
        <NavLink to="/queues/ppm" activeClassName="usa-current">
          <span>All PPMs Queue</span>
        </NavLink>
      );
    case '/queues/ppm_payment_requested':
      return (
        <NavLink to="/queues/ppm_payment_requested" activeClassName="usa-current">
          <span>Payment requested</span>
        </NavLink>
      );
    case '/queues/all':
      return (
        <NavLink to="/queues/all" activeClassName="usa-current">
          <span>All moves</span>
        </NavLink>
      );
    default:
      return (
        <NavLink to="/queues/new" activeClassName="usa-current">
          <span>New moves</span>
        </NavLink>
      );
  }
};

class MoveInfo extends Component {
  state = {
    redirectToHome: false,
    hideTooltip: true,
  };

  componentDidMount() {
    const { moveId } = this.props;
    this.props.loadMove(moveId);
    this.props.getMoveDocumentsForMove(moveId);
    this.props.loadPPMs(moveId);
  }

  componentDidUpdate(prevProps) {
    const { loadBackupContacts, loadOrders, loadMoveIsSuccess, loadServiceMember, ordersId, serviceMemberId } =
      this.props;
    if (loadMoveIsSuccess !== prevProps.loadMoveIsSuccess && loadMoveIsSuccess) {
      loadOrders(ordersId);
      loadServiceMember(serviceMemberId);
      loadBackupContacts(serviceMemberId);
    }
  }

  componentWillUnmount() {
    this.props.resetRequests();
  }

  get allAreApproved() {
    const { moveStatus, ppm } = this.props;
    const moveApproved = moveStatus === 'APPROVED';
    const ppmApproved = includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    return moveApproved && ppmApproved;
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.moveId);
  };

  approvePPM = () => {
    const approveDate = moment().format();
    this.props.approvePPM(this.props.ppm.id, approveDate);
  };

  cancelMoveAndRedirect = (cancelReason) => {
    const messageLines = [
      `Move #${this.props.move.locator} for ${this.props.serviceMember.last_name}, ${this.props.serviceMember.first_name} has been canceled`,
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
            <FontAwesomeIcon className="icon approval-ready" icon="check" />
            Move pending
          </span>
        );
      } else {
        return (
          <span className="status">
            <FontAwesomeIcon className="icon approval-waiting" icon="clock" />
            Payment Requested
          </span>
        );
      }
    } else {
      return (
        <span className="status">
          <FontAwesomeIcon className="icon approval-waiting" icon="clock" />
          In review
        </span>
      );
    }
  };

  handleToolTipHover = () => {
    // Temporarily disable due to bug: https://ustcdp3.slack.com/archives/CP4979J0G/p1575412461036700
    // this.setState({ hideTooltip: !this.state.hideTooltip });
  };

  render() {
    const { move, moveId, moveDocuments, moveStatus, orders, ppm, serviceMember, upload } = this.props;
    const showDocumentViewer = this.props.context.flags.documentViewer;
    const moveInfoComboButton = this.props.context.flags.moveInfoComboButton;
    const ordersComplete = Boolean(
      orders.orders_number && orders.orders_type_detail && orders.department_indicator && orders.tac && orders.sac,
    );
    const ppmPaymentRequested = includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    const ppmApproved = includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    const moveApproved = moveStatus === 'APPROVED';

    const moveDate = ppm.original_move_date;

    const uploadDocumentUrl = `/moves/${moveId}/documents/new`;
    const ordersUrl = `/moves/${move.id}/orders`;

    if (this.state.redirectToHome) {
      return <Redirect to="/" />;
    }

    if (!this.props.loadDependenciesHasSuccess && !this.props.loadDependenciesHasError) return <LoadingPlaceholder />;
    if (this.props.loadDependenciesHasError)
      return (
        <div className="grid-container-widescreen usa-prose">
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                Something went wrong contacting the server.
              </Alert>
            </div>
          </div>
        </div>
      );

    return (
      <div className="grid-container-widescreen usa-prose">
        <div className="grid-row grid-gap">
          <div className="grid-col-8">
            <h1>
              Move Info: {serviceMember.last_name}, {serviceMember.first_name}
            </h1>
          </div>
          <div className="grid-col-4 nav-controls">
            <ReferrerQueueLink history={this.props.history} />
          </div>
        </div>
        <div className="grid-row">
          <div className="grid-col-12">
            <ul className="move-info-header-meta">
              <li>ID# {serviceMember.edipi}&nbsp;</li>
              <li>
                {serviceMember.telephone}
                {serviceMember.phone_is_preferred && (
                  <FontAwesomeIcon className="icon icon-grey" icon="phone" flip="horizontal" />
                )}
                {serviceMember.email_is_preferred && <FontAwesomeIcon className="icon icon-grey" icon="envelope" />}
                &nbsp;
              </li>
              <li>Locator# {move.locator}&nbsp;</li>
              <li>Move date {formatDate(moveDate)}&nbsp;</li>
            </ul>
          </div>
        </div>

        <div className="grid-row grid-gap tabs">
          <div className="grid-col-9">
            <RoutedTabs startPathWith={this.props.match.url}>
              <NavTab to="/basics">
                <span className="title" data-testid="basics-tab">
                  Basics
                </span>
                <span className="status">
                  <FontAwesomeIcon className="icon" icon="play-circle" />
                  {capitalize(this.props.moveStatus)}
                </span>
              </NavTab>
              <NavTab to="/ppm">
                <span className="title" data-testid="ppm-tab">
                  PPM
                </span>
                {this.renderPPMTabStatus()}
              </NavTab>
            </RoutedTabs>

            <div className="tab-content">
              <Switch>
                <PrivateRoute
                  exact
                  path={`${this.props.match.url}`}
                  render={() => (
                    <Redirect
                      replace
                      to={{ pathname: `${this.props.match.url}/basics`, state: this.props.history.location.state }}
                    />
                  )}
                  requiredRoles={[roleTypes.PPM]}
                />
                <PrivateRoute path={`${this.props.match.path}/basics`} requiredRoles={[roleTypes.PPM]}>
                  <BasicsTabContent moveId={moveId} serviceMember={this.props.serviceMember} />
                </PrivateRoute>
                <PrivateRoute path={`${this.props.match.path}/ppm`} requiredRoles={[roleTypes.PPM]}>
                  <PPMTabContent
                    ppmPaymentRequestedFlag={this.props.context.flags.ppmPaymentRequest}
                    moveId={moveId}
                    ppmPaymentRequested={ppmPaymentRequested}
                    moveDocuments={moveDocuments}
                  />
                </PrivateRoute>
              </Switch>
            </div>
          </div>
          <div className="grid-col-3">
            <div>
              {this.props.approveMoveHasError && (
                <Alert type="warning" heading="Unable to approve">
                  Please fill out missing data
                </Alert>
              )}
              <div onMouseEnter={this.handleToolTipHover} onMouseLeave={this.handleToolTipHover}>
                <ToolTip
                  disabled={this.state.hideTooltip}
                  textStyle="tooltiptext-large"
                  toolTipText="Some information about the move is missing or contains errors. Please fix these problems before approving."
                >
                  {moveInfoComboButton && (
                    <ComboButton
                      allAreApproved={this.allAreApproved}
                      buttonText={`Approve${this.allAreApproved ? 'd' : ''}`}
                      disabled={this.allAreApproved || !ordersComplete}
                    >
                      <DropDown>
                        <DropDownItem
                          value="Approve Basics"
                          disabled={moveApproved || !ordersComplete}
                          onClick={this.approveBasics}
                        />
                        <DropDownItem
                          disabled={ppmApproved || !moveApproved || !ordersComplete}
                          onClick={this.approvePPM}
                          value="Approve PPM"
                        />
                      </DropDown>
                    </ComboButton>
                  )}
                </ToolTip>
                <ConfirmWithReasonButton
                  buttonTitle="Cancel Move"
                  reasonPrompt="Why is the move being canceled?"
                  warningPrompt="Are you sure you want to cancel the entire move?"
                  onConfirm={this.cancelMoveAndRedirect}
                  buttonDisabled={false}
                />
              </div>
            </div>
            <div className="documents">
              <h2 className="extras usa-heading">Documents</h2>
              {!upload ? (
                <p>No orders have been uploaded.</p>
              ) : (
                <div>
                  {moveApproved ? (
                    <div className="panel-field">
                      <FontAwesomeIcon style={{ color: 'green' }} className="icon" icon="check" />
                      <a
                        href={ordersUrl}
                        target={`orders-${moveId}`}
                        onClick={openLinkInNewWindow.bind(
                          this,
                          ordersUrl,
                          `orders-${moveId}`,
                          window,
                          defaultRelativeWindowSize,
                        )}
                        className="usa-link"
                      >
                        Orders ({formatDate(upload.created_at)})
                      </a>
                    </div>
                  ) : (
                    <div className="panel-field">
                      <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon="exclamation-circle" />
                      <a
                        href={ordersUrl}
                        target={`orders-${moveId}`}
                        onClick={openLinkInNewWindow.bind(
                          this,
                          ordersUrl,
                          `orders-${moveId}`,
                          window,
                          defaultRelativeWindowSize,
                        )}
                        className="usa-link"
                      >
                        Orders ({formatDate(upload.created_at)})
                      </a>
                    </div>
                  )}
                </div>
              )}
              {showDocumentViewer && (
                <DocumentList
                  detailUrlPrefix={`/moves/${moveId}/documents`}
                  moveDocuments={moveDocuments}
                  uploadDocumentUrl={uploadDocumentUrl}
                  moveId={moveId}
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
  const ppm = selectActivePPMForMove(state, moveId);
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
    ppmAdvance: selectReimbursementById(state, ppm.advance) || {},
    serviceMember,
    serviceMemberId,
    swaggerError: get(state, 'swagger.hasErrored'),
    upload: get(orders, 'uploaded_orders.uploads.0', {}),
  };
};

const mapDispatchToProps = (dispatch) =>
  bindActionCreators(
    {
      getMoveDocumentsForMove,
      approveBasics,
      approvePPM,
      cancelMove,
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

const connectedMoveInfo = withContext(connect(mapStateToProps, mapDispatchToProps)(MoveInfo));
export { connectedMoveInfo as default, ReferrerQueueLink };
