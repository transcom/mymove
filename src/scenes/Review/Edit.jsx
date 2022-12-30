import React from 'react';
import Summary from 'components/Customer/Review/Summary/Summary';
import scrollToTop from 'shared/scrollToTop';
import withRouter from 'utils/routing';

export class Edit extends React.Component {
  componentDidMount() {
    scrollToTop();
  }

  goHome = () => {
    const {
      router: { navigate },
    } = this.props;
    navigate('/');
  };

  render() {
    return (
      <div className="grid-container usa-prose">
        <div className="grid-row">
          <div className="grid-col-12">
            <a className="usa-link back-to-home" onClick={this.goHome}>
              &lt; BACK TO HOME
            </a>
            <h1 className="edit-title">Edit Move</h1>
            <p>Changes to your orders or shipments could impact your move, including the estimated PPM incentive.</p>
            <Summary />
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(Edit);
