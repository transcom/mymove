import React from 'react';
import Summary from './Summary';
import { withRouter } from 'react-router-dom';
export class Edit extends React.Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  goBack = () => {
    this.props.history.goBack();
  };
  render() {
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">
          <a className="back-to-home" onClick={this.goBack}>
            &lt;BACK TO HOME
          </a>
          <h1 className="edit-title">Edit Move</h1>
          <p>Changes to your orders or shipments could impact your move, including the estimated PPM incentive.</p>
          <Summary />
        </div>
      </div>
    );
  }
}

export default withRouter(Edit);
