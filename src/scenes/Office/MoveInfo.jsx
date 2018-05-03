import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPhone from '@fortawesome/fontawesome-free-solid/faPhone';
import faComments from '@fortawesome/fontawesome-free-solid/faComments';

import './office.css';

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
            <p>TABS GO HERE</p>
            <p>Displaying move {this.props.match.params.moveID}.</p>
          </div>
        </div>
      </div>
    );
  }
}
