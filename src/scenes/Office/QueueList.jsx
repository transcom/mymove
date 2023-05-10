import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

export default class QueueList extends Component {
  render() {
    return (
      <div>
        <h2 className="queue-list-heading">Queues</h2>
        <ul className="usa-sidenav">
          <li className="usa-sidenav__item">
            <NavLink to="/queues/new" className={({ isActive }) => (isActive ? 'usa-current' : '')}>
              <span>New moves</span>
            </NavLink>
          </li>

          <li className="usa-sidenav__item">
            <NavLink
              to="/queues/ppm_approved"
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              data-testid="ppm-payment-requests-queue"
            >
              <span>Approved</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink
              to="/queues/ppm_payment_requested"
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              data-testid="ppm-payment-requests-queue"
            >
              <span>Payment requested</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink
              to="/queues/ppm_completed"
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              data-testid="ppm-payment-requests-queue"
            >
              <span>Completed</span>
            </NavLink>
          </li>
          <li className="usa-sidenav__item">
            <NavLink
              to="/queues/all"
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              data-testid="ppm-queue"
            >
              <span>All moves</span>
            </NavLink>
          </li>
        </ul>
      </div>
    );
  }
}
