import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    return (
      <div>
        <h2 class="queue-list-heading">Queues</h2>
        <ul className="usa-sidenav-list">
          <li>
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/accepted" activeClassName="usa-current">
              <span>Accepted Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/approved" activeClassName="usa-current">
              <span>Approved Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/in_transit" activeClassName="usa-current">
              <span>In Transit Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/delivered" activeClassName="usa-current">
              <span>Delivered Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/completed" activeClassName="usa-current">
              <span>Completed Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="/queues/all" activeClassName="usa-current">
              <span>All Shipments</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
