import React, { Component } from 'react';
import ordersComplete from 'shared/images/orders-complete-gray-icon.png';
import moveIcon from 'shared/images/move-icon.png';

export class TransitionToMove extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Transition to Scheduling Move';
  }
  render() {
    return (
      <div className="usa-grid">
        <div className="lg center">
          <p> Great, we're done with your orders.</p>
          <img className="sm" src={ordersComplete} alt="profile-check" />
        </div>

        <div className="lg center">
          <p>Now, we're ready to schedule your move!</p>
          <img className="sm" src={moveIcon} alt="onto-move-orders" />
        </div>
      </div>
    );
  }
}

export default TransitionToMove;
