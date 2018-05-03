import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

import { RoutedTabs, NavTab } from 'react-router-tabs';
import { Route, Switch, Redirect } from 'react-router-dom';
import TextBoxWithEditButton from 'shared/TextBoxWithEditButton';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';
import faExclamationTriangle from '@fortawesome/fontawesome-free-solid/faExclamationTriangle';
import faPlayCircle from '@fortawesome/fontawesome-free-solid/faPlayCircle';

import './office.css';

const BasicsTabContent = () => {
  return (
    <div>
      <TextBoxWithEditButton />
    </div>
  );
};

const PPMTabContent = () => {
  return <div>PPM</div>;
};

export default class MoveInfo extends Component {
  render() {
    return (
      <div>
        <div className="usa-grid grid-wide">
          <div className="usa-width-two-thirds Todo">
            <h1>Move Info: Johnson, Casey</h1>
          </div>
          <div className="usa-width-one-third nav-controls">
            <NavLink to="/queues/new_moves" activeClassName="usa-current">
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
          <div className="usa-width-one-whole Todo">
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
                  path={`${this.props.match.url}/basics`}
                  component={BasicsTabContent}
                />
                <Route
                  path={`${this.props.match.url}/ppm`}
                  component={PPMTabContent}
                />
              </Switch>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
