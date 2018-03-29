import React, { Component } from 'react';
import { connect } from 'react-redux';
import windowSize from 'react-window-size';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import { setPendingMoveType } from './ducks';
import BigButton from 'shared/BigButton';
import trailerGray from 'shared/icon/trailer-gray.svg';
import truckGray from 'shared/icon/truck-gray.svg';
import hhgPpmCombo from 'shared/icon/hhg-ppm-combo.svg';
import './MoveType.css';

class BigButtonGroup extends Component {
  constructor() {
    super();
    this.state = {
      isHidden: true,
    };
  }
  toggleHidden() {
    this.setState({
      isHidden: !this.state.isHidden,
    });
  }
  render() {
    const isMobile = this.props.windowWidth < 481;
    const createButton = (
      value,
      description,
      title,
      icon,
      prosList,
      altTag,
      isMobile,
    ) => {
      const onButtonClick = () => {
        this.props.onMoveTypeSelected(value);
      };
      return (
        <BigButton
          value={value}
          selected={this.props.selectedOption === value}
          onClick={onButtonClick}
          className="move-type-button"
        >
          <div>
            <p className="restrict-left">{description}</p>
            <img src={icon} alt={altTag} />
            {!isMobile && <p className="font-2">{title}</p>}
            {isMobile && (
              <div
                className="collapse-btn"
                onClick={this.toggleHidden.bind(this)}
              >
                &gt; &nbsp; Pros and Cons:
              </div>
            )}
            {(!isMobile || !this.state.isHidden) && (
              <div>
                {Object.keys(prosList || {}).map(function(key) {
                  const pros = prosList[key];
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
            )}
          </div>
        </BigButton>
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
      isMobile,
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
      isMobile,
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
      isMobile,
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
  windowWidth: PropTypes.number,
  selectedOption: PropTypes.string,
  onMoveTypeSelected: PropTypes.func,
};

const BigButtonGroupWithSize = windowSize(BigButtonGroup);

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
        <h2> Select a Move Type</h2>
        <BigButtonGroupWithSize
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
    selected_move_type: PropTypes.string,
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
