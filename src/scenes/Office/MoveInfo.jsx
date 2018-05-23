import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Switch, Redirect, Link } from 'react-router-dom';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'shared/User/PrivateRoute';
import Alert from 'shared/Alert'; // eslint-disable-line
import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';
import PaymentsPanel from './PaymentsPanel';
import { loadMoveDependencies, approveBasics } from './ducks.js';
import { formatDate } from './helpers';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

import './office.css';

const BasicsTabContent = props => {
  return (
    <React.Fragment>
      <OrdersPanel title="Orders" moveId={props.match.params.moveId} />
      <CustomerInfoPanel
        title="Customer Info"
        moveId={props.match.params.moveId}
      />
      <BackupInfoPanel title="Backup Info" moveId={props.match.params.moveId} />
      <AccountingPanel title="Accounting" moveId={props.match.params.moveId} />
    </React.Fragment>
  );
};

const PPMTabContent = props => {
  return (
    <React.Fragment>
      <PaymentsPanel title="Payments" moveId={props.match.params.moveId} />
    </React.Fragment>
  );
};

class MoveInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.match.params.moveId);
  };

  render() {
    // TODO: If the following vars are not used to load data, remove them.
    const officeMove = get(this.props, 'officeMove', {});
    // const officeOrders = this.props.officeOrders || {};
    const officeServiceMember = get(this.props, 'officeServiceMember', {});
    // const officeBackupContacts = this.props.officeBackupContacts || []
    // Todo: Change once more than 1 PPM will be loaded at one time
    const officePPM = get(this.props, 'officePPMs[0]');

    let upload = get(this.props, 'officeOrders.uploaded_orders.uploads.0'); // there can be only one

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
              Move Info: {officeServiceMember.last_name},{' '}
              {officeServiceMember.first_name}
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
              <li>ID# {officeServiceMember.id}</li>
              <li>
                {officeServiceMember.telephone}
                {officeServiceMember.phone_is_preferred && (
                  <FontAwesomeIcon
                    className="icon"
                    icon={faPhone}
                    flip="horizontal"
                  />
                )}
                {officeServiceMember.text_message_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faComments} />
                )}
                {officeServiceMember.email_is_preferred && (
                  <FontAwesomeIcon className="icon" icon={faEmail} />
                )}
              </li>
              <li>Locator# {officeMove.locator}</li>
              <li className="Todo">KKFA to HAFC</li>
              <li>
                Move date {formatDate(get(officePPM, 'planned_move_date'))}
              </li>
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
                  {capitalize(officeMove.status)}
                </span>
              </NavTab>
              <NavTab to="/ppm">
                <span className="title">PPM</span>
                {officePPM.status === 'APPROVED' ? (
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon approval-ready"
                      icon={faCheck}
                    />
                    Move pending
                  </span>
                ) : (
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon approval-waiting"
                      icon={faClock}
                    />
                    In review
                  </span>
                )}
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
                disabled={officeMove.status === 'APPROVED'}
              >
                Approve Basics
              </button>
              <button>Troubleshoot</button>
              <button>Cancel Move</button>
            </div>
            <div className="documents">
              <h2 className="usa-heading">
                Documents
                <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
              </h2>
              {!upload ? (
                <p>No orders have been uploaded.</p>
              ) : (
                <div className="document">
                  <FontAwesomeIcon
                    style={{ color: 'red' }}
                    className="icon"
                    icon={faExclamationCircle}
                  />
                  <Link to={`/moves/${officeMove.id}/orders`} target="_blank">
                    Orders ({formatDate(upload.created_at)})
                  </Link>
                </div>
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
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  officeMove: state.office.officeMove,
  officeOrders: state.office.officeOrders,
  officeServiceMember: state.office.officeServiceMember,
  officeBackupContacts: state.office.officeBackupContacts,
  officePPMs: state.office.officePPMs,
  loadDependenciesHasSuccess: state.office.loadDependenciesHasSuccess,
  loadDependenciesHasError: state.office.loadDependenciesHasError,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMoveDependencies, approveBasics }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(MoveInfo);
