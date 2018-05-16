import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';
import moment from 'moment';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Route, Switch, Redirect } from 'react-router-dom';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import PrivateRoute from 'shared/User/PrivateRoute';
import Alert from 'shared/Alert'; // eslint-disable-line
import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';
import {
  loadMoveDependencies,
  loadAccounting,
  approveBasics,
} from './ducks.js';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faEmail from '@fortawesome/fontawesome-free-solid/faEnvelope';
import faExclamationTriangle from '@fortawesome/fontawesome-free-solid/faExclamationTriangle';
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

const PPMTabContent = () => {
  return <div>PPM</div>;
};

class MoveInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
    this.props.loadAccounting(this.props.match.params.moveId);
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.match.params.moveId);
  };

  render() {
    // TODO: If the following vars are not used to load data, remove them.
    const officeMove = this.props.officeMove || {};
    // const officeOrders = this.props.officeOrders || {};
    const officeServiceMember = this.props.officeServiceMember || {};
    // const officeBackupContacts = this.props.officeBackupContacts || []
    const officePPMs = this.props.officePPMs || [];

    let uploads;
    if (this.props.officeOrders) {
      uploads = this.props.officeOrders.uploaded_orders.uploads;
    } else {
      uploads = [];
    }

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
              <li className="Todo">Locator# {officeMove.locator}</li>
              <li>KKFA to HAFC</li>
              <li>
                Requested Pickup {get(officePPMs, '[0].planned_move_date')}
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
                <span className="status">
                  <FontAwesomeIcon
                    className="icon"
                    icon={faExclamationTriangle}
                  />
                  Status Goes Here
                </span>
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
              {uploads.length === 0 && <p>No orders have been uploaded.</p>}
              {uploads.map(upload => {
                return (
                  <div key={upload.url} className="document">
                    <FontAwesomeIcon
                      style={{ color: 'red' }}
                      className="icon"
                      icon={faExclamationCircle}
                    />
                    <a href={upload.url} target="_blank">
                      Orders ({moment(upload.created_at).format('D-MMM-YY')})
                    </a>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      </div>
    );
  }
}

MoveInfo.propTypes = {
  loadMoveDependencies: PropTypes.func.isRequired,
  loadAccounting: PropTypes.func.isRequired,
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
  bindActionCreators(
    { loadMoveDependencies, loadAccounting, approveBasics },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(MoveInfo);
