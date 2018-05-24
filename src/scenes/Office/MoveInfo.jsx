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
import PPMEstimatesPanel from './PPMEstimatesPanel';
import OrdersPanel from './OrdersPanel';
import { loadMoveDependencies, approveBasics, approvePPM } from './ducks.js';
import { formatDate } from './helpers';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
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
  return <PPMEstimatesPanel />;
};

class MoveInfo extends Component {
  componentDidMount() {
    this.props.loadMoveDependencies(this.props.match.params.moveId);
  }

  approveBasics = () => {
    this.props.approveBasics(this.props.match.params.moveId);
  };

  approvePPM = () => {
    this.props.approvePPM(this.props.officeMove.id, this.props.officePPM.id);
  };

  render() {
    const move = this.props.officeMove;
    const serviceMember = this.props.officeServiceMember;
    const ppm = this.props.officePPM;

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
              <li>ID# {serviceMember.id}</li>
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
              <li className="Todo">KKFA to HAFC</li>
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
                <div>
                  {move.status === 'APPROVED' ? (
                    <div className="document">
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
                    <div className="document">
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
  swaggerError: get(state, 'swagger.hasErrored'),
  officeMove: get(state, 'office.officeMove', {}),
  officeOrders: get(state, 'office.officeOrders', {}),
  officeServiceMember: get(state, 'office.officeServiceMember', {}),
  officeBackupContacts: get(state, 'office.officeBackupContacts', []),
  officePPM: get(state, 'office.officePPMs.0', {}),
  loadDependenciesHasSuccess: get(state, 'office.loadDependenciesHasSuccess'),
  loadDependenciesHasError: get(state, 'office.loadDependenciesHasError'),
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    { loadMoveDependencies, approveBasics, approvePPM },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(MoveInfo);
