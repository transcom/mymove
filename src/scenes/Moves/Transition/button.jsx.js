import React, { Component } from 'react';
import PropTypes from 'prop-types';

class MovingPage extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      selectedButton: '',
    };
  }

  render() {
    var selectedButton = this.state.selectedButton;

    var createButton = (label, firstLine, secondLine, icon) => {
      var onButtonClick = () => {
        this.setState({ selectedButton: label });
        // This is how data is usually passed up from children to their parents...
        this.props.onMoveTypeSelected(label);
        // but the problem is that it doesn't scale if if you want to pass it to your
        // parent's parent, or their parent, or their parent... i.e. great-great-...parent.
        // Solution 1: Redux is how people try to solve this. basically a big global bag
        // of shared state.
        // Solution 2: there are lighter weight solutions than Redux. Redux can be
        // overkill for a lot of things.
      };
      return (
        <MovingButton
          firstLine={firstLine}
          secondLine={secondLine}
          icon={icon}
          selected={this.state.selectedButton == label}
          onButtonClick={onButtonClick}
        />
      );
    };

    var small = createButton(
      'small',
      'this is a small move',
      '100-200 lbs',
      'http://www.crystalinks.com/kangaroo.jpg',
    );
    var medium = createButton(
      'medium',
      'this is a medium move',
      '100-200 lbs',
      'http://www.crystalinks.com/kangaroo.jpg',
    );
    var large = createButton(
      'large',
      'this is a large move',
      '100-200 lbs',
      'http://www.crystalinks.com/kangaroo.jpg',
    );

    return (
      <div>
        {small}
        {medium}
        {large}
      </div>
    );
  }
}

MovingPage.PropTypes = {
  onMoveTypeSelected: React.PropTypes.func,
};

class MovingButton extends React.Component {
  render() {
    let className = 'button';
    if (this.props.selected) {
      className += ' selected';
    }
    return (
      <div className={className} onClick={this.props.onButtonClick}>
        <div>{this.props.firstLine}</div>
        <div>{this.props.secondLine}</div>
        <img className="icon" src={this.props.icon} />
      </div>
    );
  }
}

MovingButton.propTypes = {
  firstLine: React.PropTypes.string.isRequired,
  secondLine: React.PropTypes.string.isRequired,
  icon: React.PropTypes.string.isRequired,
  selected: React.PropTypes.bool,
  onButtonClick: React.PropTypes.func,
};

/*
 * Render the above component into the div#app
 */
ReactDOM.render(<MovingPage />, document.getElementById('app'));
