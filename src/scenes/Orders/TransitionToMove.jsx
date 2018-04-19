import React, { Component } from 'react';
import ordersComplete from 'shared/images/orders_check.png';
import moveIcon from 'shared/images/move_icon.png';

export class TransitionToMove extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Transition to Scheduling Move';
  }
  render() {
    return (
      <div className="usa-grid">
        <div className="lg center">
          <p> Great, we're done with your orders.</p>
          <img className="sm Todo" src={ordersComplete} alt="profile-check" />
        </div>

        <div className="lg center">
          <p>Now, we're ready to schedule your move!</p>
          <img className="sm Todo" src={moveIcon} alt="onto-move-orders" />
        </div>

        <ul className="Todo">
          <li>
            Replace screenshots of icons with svgs once we get assets from
            design.
          </li>
        </ul>
      </div>
    );
  }
}

export default TransitionToMove;
