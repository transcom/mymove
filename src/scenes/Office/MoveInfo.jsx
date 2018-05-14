import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';
import PropTypes from 'prop-types';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Route, Switch, Redirect } from 'react-router-dom';

import AccountingPanel from './AccountingPanel';
import BackupInfoPanel from './BackupInfoPanel';
import CustomerInfoPanel from './CustomerInfoPanel';
import OrdersPanel from './OrdersPanel';

import { loadMove, loadAccounting } from './ducks.js';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faExclamationTriangle from '@fortawesome/fontawesome-free-solid/faExclamationTriangle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';

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
    this.props.loadMove(this.props.match.params.moveId);
    this.props.loadAccounting(this.props.match.params.moveId);
  }

  render() {
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds Todo">
            <h1>Move Info: Johnson, Casey</h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves Queue</span>
            </NavLink>
          </div>
        </div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-one-whole Todo">
            <ul className="move-info-header-meta">
              <li>ID# 3938593893</li>
              <li>
                (303) 936-8181
                <FontAwesomeIcon
                  className="icon"
                  icon={faPhone}
                  flip="horizontal"
                />
                <FontAwesomeIcon className="icon" icon={faComments} />
              </li>
              <li>Locator# ABC89</li>
              <li>KKFA to HAFC</li>
              <li>Requested Pickup 5/10/18</li>
            </ul>
          </div>
        </div>

        <div className="usa-grid grid-wide tabs">
          <div className="usa-width-three-fourths">
            <p>Displaying move {this.props.match.params.moveID}.</p>

            <RoutedTabs startPathWith={this.props.match.url}>
              <NavTab to="/basics">
                <span className="title">Basics</span>
                <span className="status">
                  <FontAwesomeIcon className="icon" icon={faPlayCircle} />
                  Status Goes Here
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
                <Route
                  exact
                  path={`${this.props.match.url}`}
                  render={() => (
                    <Redirect replace to={`${this.props.match.url}/basics`} />
                  )}
                />
                <Route
                  path={`${this.props.match.path}/basics`}
                  component={BasicsTabContent}
                />
                <Route
                  path={`${this.props.match.path}/ppm`}
                  component={PPMTabContent}
                />
              </Switch>
            </div>
          </div>
          <div className="usa-width-one-fourths">
            <div>
              <button>Approve Basics</button>
              <button>Troubleshoot</button>
              <button>Cancel Move</button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

MoveInfo.propTypes = {
  loadMove: PropTypes.func.isRequired,
  loadAccounting: PropTypes.func.isRequired,
};

const mapStateToProps = state => ({
  swaggerError: state.swagger.hasErrored,
  officeMove: state.office.officeMove,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadMove, loadAccounting }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(MoveInfo);
