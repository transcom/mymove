import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

class Wizard extends Component {
  constructor(props) {
    super(props);
    this.nextPage = this.nextPage.bind(this);
    this.previousPage = this.previousPage.bind(this);
    this.state = {
      currentPageIndex: 0,
    };
  }
  nextPage() {
    this.setState({ currentPageIndex: this.state.currentPageIndex + 1 });
  }

  previousPage() {
    this.setState({ currentPageIndex: this.state.currentPageIndex - 1 });
  }

  render() {
    const { children } = this.props;
    const { currentPageIndex } = this.state;
    const lastPageIndex = React.Children.count(children) - 1;
    const isFirstPage = currentPageIndex === 0;
    const isLastPage = currentPageIndex === lastPageIndex;
    const getCurrentPage = () => {
      return React.Children.map(children, (child, i) => {
        if (i !== currentPageIndex) return;
        return child;
      });
    };
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">{getCurrentPage()}</div>
        <div className="usa-width-one-third">
          {!isFirstPage && (
            <button
              className={classnames({ 'usa-button-secondary': !isLastPage })}
              onClick={this.previousPage}
            >
              Prev
            </button>
          )}
        </div>
        <div className="usa-width-one-third" />
        <div className="usa-width-one-third">
          {!isLastPage && <button onClick={this.nextPage}>Next</button>}
        </div>
      </div>
    );
  }
}

Wizard.propTypes = {
  children: PropTypes.node,
};

export default Wizard;
