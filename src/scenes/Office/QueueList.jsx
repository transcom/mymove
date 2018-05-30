import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    return (
      <div>
        <h2>Queues</h2>
        <ul className="usa-sidenav-list">
          <li>
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/ppm" activeClassName="usa-current">
              <span>PPMs</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
