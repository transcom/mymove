import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { setPendingMoveType } from './ducks';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import hhgPpmCombo from 'shared/icon/hhg-ppm-combo.svg';
import './MoveType.css';

export class BigButton extends Component {
  render() {
    let className = 'move-type-button';
    if (this.props.selected) {
      className += ' selected';
    }
    // const imgs = this.props.icons.map(icon => (
    //   <img src={icon} alt={this.props.altTag} key={icon} />
    // ));
    return (
      <div className={className} onClick={this.props.onButtonClick}>
        <p className="restrict-left">{this.props.description}</p>
        <img src={this.props.icon} alt={this.props.altTag} />
        <p className="font-2">{this.props.title}</p>
        {Object.keys(this.props.pros || {}).map(function(key) {
          var pros = this.props.pros[key];
          return (
            <div key={key.toString()}>
              <p>{key}</p>
              <ul className="font-3">
                {pros.map(item => <li key={item}>{item}</li>)}
              </ul>
            </div>
          );
        }, this)}
        <p className="move-type-button-more-info">
          <a href="about:blank">more information</a>
        </p>
      </div>
    );
  }
}

BigButton.propTypes = {
  description: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  icon: PropTypes.string.isRequired,
  altTag: PropTypes.string.isRequired,
  pros: PropTypes.object.isRequired,
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
    var combo = createButton(
      'COMBO',
      'Government moves the big stuff and you move the rest',
      'HHG Move with Partial PPM',
      hhgPpmCombo,
      {
        'Pros:': [
          'The government can arrange a mover to handle big stuff.',
          'Potential for you to earn a little $ by moving some items yourself.',
          'Protect valuable or sentimental items by moving them with you.',
        ],
        'Cons:': [
          'More things to keep track of.',
          'More prep to separate what you move from what the gov moves.',
          'Have the overhead work of PPM for potentially minimal incentive.',
        ],
      },
      'hhg-ppm-combo',
    );
    var ppm = createButton(
      'PPM',
      'You move 100% of your household goods',
      'Personally Procured Move (PPM)',
      trailerGray,
      {
        'Pros:': [
          'You choose how your stuff is transported.',
          'Potential to earn a small amount of $.',
          'Flexible moving dates (during the week, weekend, or across multiple trips/dates).',
          'Can still hire moving company or use a pod.',
        ],
        'Cons:': [
          'You have to arrange everything.',
          'More work: packing, weighing, transporting, submitting paperwork.',
          'You can only submit claims for things that are not your fault.',
        ],
      },
      'trailer-gray',
    );
    var hhg = createButton(
      'HHG',
      'Government handles 100% of your move',
      'Household Goods Move (HHG)',
      truckGray,
      {
        'Pros:': [
          'The government arranges moving companies to pack & transport your stuff.',
          'Less hassle',
          'The claims process is available to you if anything becomes damaged/broken).',
        ],
        'Cons:': [
          'Limited availability.',
          'Can only move on weekdays.',
          'Your stuff is placed in storage if you cannot meet the truck at the destination.',
          'You may not like your moving company.',
        ],
      },
      'truck-gray',
    );

    return (
      <div className="move-type-content">
        <div className="usa-width-one-third">{combo}</div>
        <div className="usa-width-one-third">{ppm}</div>
        <div className="usa-width-one-third">{hhg}</div>
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
    const selectedOption = pendingMoveType;
    return (
      <div className="usa-grid-full">
        <h2> Select a Move Type</h2>
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
