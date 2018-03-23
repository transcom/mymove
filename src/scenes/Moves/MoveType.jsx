import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { setPendingMoveType } from './ducks';
import carGray from 'shared/icon/car-gray.svg';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import './MoveType.css';

export class BigButton extends Component {
  render() {
    let className = 'ppm-size-button';
    if (this.props.selected) {
      className += ' selected';
    }
    return (
      <div className={className} onClick={this.props.onButtonClick}>
        <p>{this.props.description}</p>
        <img className="icon" src={this.props.icon} alt={this.props.altTag} />
        <p>{this.props.title}</p>
        // <div>{this.props.pros}</div>
        <div>
          <p> Pros </p>
          <p>1 sldfkjsl</p>
        </div>
      </div>
    );
  }
}

BigButton.propTypes = {
  description: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  icon: PropTypes.string.isRequired,
  altTag: PropTypes.string.isRequired,
  pros: PropTypes.string.isRequired,
  selected: PropTypes.bool,
  onButtonClick: PropTypes.func,
};

export class BigButtonGroup extends Component {
  render() {
    var createButton = (value, description, title, icon, pros, altTag) => {
      var onButtonClick = () => {
        this.props.onMoveTypeSelected(value);
      };
      return (
        <BigButton
          value={value}
          description={description}
          title={title}
          icon={icon}
          pros={pros}
          altTag={altTag}
          selected={this.props.selectedOption === value}
          onButtonClick={onButtonClick}
        />
      );
    };
    var hhg = createButton(
      'HHG',
      'Government moves the big stuff and you move the rest',
      'HHG Move with Partial PPM',
      trailerGray,
      'Pros: The government can arrange a mover to handle big stuff. Potential for you to earn a little $ by moving some items yourself. Protect valuable or sentimental items by moving them with you.',
      'trailer-gray',
    );
    var ppm = createButton(
      'PPM',
      'A trailer full of household goods?',
      '(approx 400 - 1,200 lbs)',
      trailerGray,
      'pros: ppm',
      'trailer-gray',
    );
    var combo = createButton(
      'COMBO',
      'A moving truck that you rent yourself?',
      '(approx 1,000 - 5,000 lbs)',
      truckGray,
      'pros: hhg',
      'truck-gray',
    );

    return (
      <div>
        <div className="usa-width-one-third">{hhg}</div>
        <div className="usa-width-one-third">{ppm}</div>
        <div className="usa-width-one-third">{combo}</div>
      </div>
    );
  }
}
BigButtonGroup.propTypes = {
  selectedOption: PropTypes.string,
  onMoveTypeSelected: PropTypes.func,
};

export class MoveType extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Move Type Selection';
  }

  onMoveTypeSelected = value => {
    this.props.setPendingMoveType(value);
  };
  render() {
    const { pendingMoveType, currentMove } = this.props;
    const selectedOption =
      pendingMoveType || (currentMove && currentMove.selected_move_type);
    return (
      <div className="usa-grid-full">
        <h3> Select a Move Type</h3>
        <BigButtonGroup
          selectedOption={selectedOption}
          onMoveTypeSelected={this.onMoveTypeSelected}
        />
      </div>
    );
  }
}

MoveType.propTypes = {
  pendingMoveType: PropTypes.string,
  currentMove: PropTypes.shape({
    id: PropTypes.string,
    size: PropTypes.string,
  }),
  setPendingMoveType: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return state.submittedMoves;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ setPendingMoveType }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(MoveType);
