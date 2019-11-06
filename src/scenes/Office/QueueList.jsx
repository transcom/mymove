import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    return (
      <div>
        <h2 className="queue-list-heading">Queues</h2>
        <ul className="usa-sidenav">
          <li className="usa-sidenav__item">
            <NavLink to="/queues/new" activeClassName="usa-current">
              <span>New moves</span>
            </NavLink>
          </li>

          <li className="usa-sidenav__item">
            <NavLink to="/queues/ppm_approved" activeClassName="usa-current" data-cy="ppm-payment-requests-queue">
              <span>Approved</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink
              to="/queues/ppm_payment_requested"
              activeClassName="usa-current"
              data-cy="ppm-payment-requests-queue"
            >
              <span>Payment requested</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink to="/queues/ppm_completed" activeClassName="usa-current" data-cy="ppm-payment-requests-queue">
              <span>Completed</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink to="/queues/all" activeClassName="usa-current" data-cy="ppm-queue">
              <span>All moves</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
