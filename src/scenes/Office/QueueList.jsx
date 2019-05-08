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
            <NavLink to="/queues/ppm" activeClassName="usa-current" data-cy="ppm-queue">
              <span>PPMs</span>
            </NavLink>
          </li>
          <li>
            <NavLink to="#hhgshipments" activeClassName="usa-current">
              <span>HHG shipments:</span>
            </NavLink>
            <ul class="usa-sidenav-sub_list">
              <li>
                <NavLink to="/queues/hhg_accepted" activeClassName="usa-current">
                  <span>Accepted</span>
                </NavLink>
              </li>
              <li>
                <NavLink to="/queues/hhg_delivered" activeClassName="usa-current">
                  <span>Delivered</span>
                </NavLink>
              </li>
              <li>
                <NavLink to="/queues/hhg_completed" activeClassName="usa-current">
                  <span>Completed</span>
                </NavLink>
              </li>
            </ul>
          </li>

          <li>
            <NavLink to="/queues/all" activeClassName="usa-current">
              <span>All Moves</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
