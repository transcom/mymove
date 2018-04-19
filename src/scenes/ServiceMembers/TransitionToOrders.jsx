import React, { Component } from 'react';
import profileComplete from 'shared/images/profile_check.png';
import moveOrders from 'shared/images/man_with_clipboard.png';
import 'scenes/ServiceMembers/ServiceMembers.css';

export class TransitionToOrders extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Transition to Station Orders';
  }
  render() {
    return (
      <div className="usa-grid">
        <div className="lg center">
          <p> OK, your profile's complete!</p>
          <img className="sm Todo" src={profileComplete} alt="profile-check" />
        </div>

        <div className="lg center">
          <p>Now, we need to take a look at your move orders.</p>
          <img className="sm Todo" src={moveOrders} alt="onto-move-orders" />
        </div>
      </div>
    );
  }
}

export default TransitionToOrders;
