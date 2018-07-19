import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { NavLink, Switch, Redirect, Link } from 'react-router-dom';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'shared/User/PrivateRoute';
import Alert from 'shared/Alert'; // eslint-disable-line
import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';
import PaymentsPanel from './Ppm/PaymentsPanel';
import PPMEstimatesPanel from './Ppm/PPMEstimatesPanel';
import StorageReimbursementCalculator from './Ppm/StorageReimbursementCalculator';
import IncentiveCalculator from './Ppm/IncentiveCalculator';
import DocumentList from 'scenes/Office/DocumentViewer/DocumentList';
import { withContext } from 'shared/AppContext';

import {
  loadMoveDependencies,
  approveBasics,
  approvePPM,
  cancelMove,
} from './ducks';
import { formatDate } from 'shared/formatters';
import {
  selectAllDocumentsForMove,
  getMoveDocumentsForMove,
} from 'shared/Entities/modules/moveDocuments';

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
    <div className="basics">
      <OrdersPanel title="Orders" moveId={props.match.params.moveId} />
      <CustomerInfoPanel
        title="Customer Info"
        moveId={props.match.params.moveId}
      />
      <BackupInfoPanel title="Backup Info" moveId={props.match.params.moveId} />
      <AccountingPanel title="Accounting" moveId={props.match.params.moveId} />
    </div>
  );
};

const PPMTabContent = props => {
  return (
    <div className="basics">
      <PaymentsPanel title="Payments" moveId={props.match.params.moveId} />
      <IncentiveCalculator />
      <StorageReimbursementCalculator />
      <PPMEstimatesPanel title="Estimates" moveId={props.match.params.moveId} />
    </div>
  );
};

class CancelPanel extends Component {
  state = { displayState: 'Button', cancelReason: '' };

  setConfirmState = () => {
    this.setState({ displayState: 'Confirm' });
  };

  setCancelState = () => {
    if (this.state.cancelReason !== '') {
      this.setState({ displayState: 'Cancel' });
    }
  };

  setButtonState = () => {
    this.setState({ displayState: 'Button' });
  };

  handleChange = event => {
    this.setState({ cancelReason: event.target.value });
  };

  cancelMove = event => {
    event.preventDefault();
    this.props.cancelMove(this.state.cancelReason);
    this.setState({ displayState: 'Redirect' });
  };

  render() {
    if (this.state.displayState === 'Cancel') {
      return (
        <div className="cancel-panel">
          <h2 className="extras usa-heading">Cancel Move</h2>
          <div className="extras content">
            <Alert type="warning" heading="Cancelation Warning">
              Are you sure you want to cancel the entire move?
            </Alert>
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>No, never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.cancelMove}>Yes, cancel move</button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'Confirm') {
      return (
        <div className="cancel-panel">
          <h2 className="extras usa-heading">Cancel Move</h2>
          <div className="extras content">
            Why is the move being canceled?
            <textarea required onChange={this.handleChange} />
            <div className="usa-grid">
              <div className="usa-width-one-whole extras options">
                <a onClick={this.setButtonState}>Never mind</a>
              </div>
              <div className="usa-width-one-whole extras options">
                <button onClick={this.setCancelState}>
                  Cancel entire move
                </button>
              </div>
            </div>
          </div>
        </div>
      );
    } else if (this.state.displayState === 'Button') {
      return (
        <button className="usa-button-secondary" onClick={this.setConfirmState}>
          Cancel Move
        </button>
      );
    } else if (this.state.displayState === 'Redirect') {
      return <Redirect to="/" />;
    }
  }
}

class MoveInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
    this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.match.params.moveId);
  };

  approvePPM = () => {
    this.props.approvePPM(this.props.officeMove.id, this.props.officePPM.id);
  };

  cancelMove = cancelReason => {
    this.props.cancelMove(this.props.officeMove.id, cancelReason);
  };

  renderPPMTabStatus = () => {
    if (this.props.officePPM.status === 'APPROVED') {
      if (
        this.props.ppmAdvance.status === 'APPROVED' ||
        !this.props.ppmAdvance.status
      ) {
        return (
          <span className="status">
            <FontAwesomeIcon className="icon approval-ready" icon={faCheck} />Move
            pending
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
    const move = this.props.officeMove;
    const serviceMember = this.props.officeServiceMember;
    const ppm = this.props.officePPM;
    const { moveDocuments } = this.props;
    const showDocumentViewer = this.props.context.flags.documentViewer;
    let upload = get(this.props, 'officeOrders.uploaded_orders.uploads.0'); // there can be only one
    let check = <FontAwesomeIcon className="icon" icon={faCheck} />;

    if (
      !this.props.loadDependenciesHasSuccess &&
      !this.props.loadDependenciesHasError
    )
      return <LoadingPlaceholder />;
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
              <li>ID# {serviceMember.edipi}</li>
              <li>
                {serviceMember.telephone}
                {serviceMember.phone_is_preferred && (
                  <FontAwesomeIcon
                    className="icon"
                    icon={faPhone}
                    flip="horizontal"
                  />
                )}
                {serviceMember.text_message_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faComments} />
                )}
                {serviceMember.email_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faEmail} />
                )}
              </li>
              <li>Locator# {move.locator}</li>
              {/*<li className="Todo">KKFA to HAFC</li>*/}
              <li>Move date {formatDate(ppm.planned_move_date)}</li>
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
              <NavTab to="/ppm">
                <span className="title">PPM</span>
                {this.renderPPMTabStatus()}
              </NavTab>
            </RoutedTabs>

            <div className="tab-content">
              <Switch>
                <PrivateRoute
                  exact
                  path={`${this.props.match.url}`}
                  render={() => (
                    <Redirect replace to={`${this.props.match.url}/basics`} />
                  )}
                />
                <PrivateRoute
                  path={`${this.props.match.path}/basics`}
                  component={BasicsTabContent}
                />
                <PrivateRoute
                  path={`${this.props.match.path}/ppm`}
                  component={PPMTabContent}
                />
              </Switch>
            </div>
          </div>
          <div className="usa-width-one-fourth">
            <div>
              <button
                onClick={this.approveBasics}
                disabled={move.status === 'APPROVED'}
                style={{
                  backgroundColor: move.status === 'APPROVED' && 'green',
                }}
              >
                Approve Basics
                {move.status === 'APPROVED' && check}
              </button>
              <button
                onClick={this.approvePPM}
                disabled={
                  ppm.status === 'APPROVED' || move.status !== 'APPROVED'
                }
                style={{
                  backgroundColor: ppm.status === 'APPROVED' && 'green',
                }}
              >
                Approve PPM
                {ppm.status === 'APPROVED' && check}
              </button>
              <CancelPanel cancelMove={this.cancelMove} />
              {/* Disabling until features implemented
              <button>Troubleshoot</button>
              */}
            </div>
            <div className="documents">
              <h2 className="extras usa-heading">
                Documents
                {!showDocumentViewer && (
                  <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
                )}
                {showDocumentViewer && (
                  <Link to={`/moves/${move.id}/documents`} target="_blank">
                    <FontAwesomeIcon
                      className="icon"
                      icon={faExternalLinkAlt}
                    />
                  </Link>
                )}
              </h2>
              {!upload ? (
                <p>No orders have been uploaded.</p>
              ) : (
                <div>
                  {move.status === 'APPROVED' ? (
                    <div className="panel-field">
                      <FontAwesomeIcon
                        style={{ color: 'green' }}
                        className="icon"
                        icon={faCheck}
                      />
                      <Link to={`/moves/${move.id}/orders`} target="_blank">
                        Orders ({formatDate(upload.created_at)})
                      </Link>
                    </div>
                  ) : (
                    <div className="panel-field">
                      <FontAwesomeIcon
                        style={{ color: 'red' }}
                        className="icon"
                        icon={faExclamationCircle}
                      />
                      <Link to={`/moves/${move.id}/orders`} target="_blank">
                        Orders ({formatDate(upload.created_at)})
                      </Link>
                    </div>
                  )}
                </div>
              )}
              {showDocumentViewer && (
                <DocumentList
                  moveDocuments={moveDocuments}
                  moveId={this.props.match.params.moveId}
                />
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

MoveInfo.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
  context: PropTypes.shape({
    flags: PropTypes.shape({ documentViewer: PropTypes.bool }).isRequired,
  }).isRequired,
};

const mapStateToProps = state => ({
  swaggerError: get(state, 'swagger.hasErrored'),
  officeMove: get(state, 'office.officeMove', {}),
  officeOrders: get(state, 'office.officeOrders', {}),
  officeServiceMember: get(state, 'office.officeServiceMember', {}),
  officeBackupContacts: get(state, 'office.officeBackupContacts', []),
  officePPM: get(state, 'office.officePPMs.0', {}),
  ppmAdvance: get(state, 'office.officePPMs.0.advance', {}),
  moveDocuments: selectAllDocumentsForMove(
    state,
    get(state, 'office.officeMove.id', ''),
  ),
  loadDependenciesHasSuccess: get(state, 'office.loadDependenciesHasSuccess'),
  loadDependenciesHasError: get(state, 'office.loadDependenciesHasError'),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      loadMoveDependencies,
      getMoveDocumentsForMove,
      approveBasics,
      approvePPM,
      cancelMove,
    },
    dispatch,
  );

export default withContext(
  connect(mapStateToProps, mapDispatchToProps)(MoveInfo),
);
