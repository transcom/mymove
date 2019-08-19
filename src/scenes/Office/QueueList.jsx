import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    // Used for the isActive attribute in NavLink.
    // True, if any params are included in the location path. False, otherwise.
    const isActive = (...args) => (match, location) => {
      return args.some(element => location.pathname.includes(element));
    };

    return (
      <div>
        <h2 className="queue-list-heading">Queues</h2>
        <ul className="usa-sidenav-list">
          <li>
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New Moves/Shipments</span>
            </NavLink>
          </li>
          <li>
            <NavLink
              to="#ppmshipments"
              activeClassName="usa-current"
              isActive={isActive('ppm', 'ppm_payment_requested')}
            >
              <span>PPM shipments:</span>
            </NavLink>
            <ul className="usa-sidenav-sub_list">
              <li>
                <NavLink to="/queues/ppm" activeClassName="usa-current" data-cy="ppm-queue">
                  <span>All</span>
                </NavLink>
              </li>
              <li>
                <NavLink
                  to="/queues/ppm_payment_requested"
                  activeClassName="usa-current"
                  data-cy="ppm-payment-requests-queue"
                >
                  <span>Payment requests</span>
                </NavLink>
              </li>
            </ul>
          </li>
          <li>
            <NavLink
              to="#hhgshipments"
              activeClassName="usa-current"
              isActive={isActive('hhg_active', 'hhg_delivered')}
            >
              <span>HHG shipments:</span>
            </NavLink>
            <ul className="usa-sidenav-sub_list">
              <li>
                <NavLink to="/queues/hhg_active" activeClassName="usa-current">
                  <span>Active</span>
                </NavLink>
              </li>
              <li>
                <NavLink to="/queues/hhg_delivered" activeClassName="usa-current">
                  <span>Delivered</span>
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
