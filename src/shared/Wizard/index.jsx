import React, { Fragment, Component } from 'react';
import PropTypes from 'prop-types';

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
    const lastPageIndex = React.Children.count(this.props.children) - 1;
    const { currentPageIndex } = this.state;
    //  const CurrentPage = children[currentPageIndex];
    return (
      <Fragment>
        {React.Children.map(children, (child, i) => {
          if (i !== currentPageIndex) return;
          return child;
        })}
        {currentPageIndex > 0 && (
          <button onClick={this.previousPage}>Prev</button>
        )}
        {currentPageIndex < lastPageIndex && (
          <button onClick={this.nextPage}>Next</button>
        )}
      </Fragment>
    );
  }
}

Wizard.propTypes = {
  children: PropTypes.node,
};

export default Wizard;
