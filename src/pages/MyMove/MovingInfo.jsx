import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { loadEntitlementsFromState } from 'shared/entitlements';

export const MovingInfo = (props) => (
  <div className="usa-grid">
    <div className="grid-row">
      <div className="grid-col">
        <h1 className="sm-heading">Figure out your move details</h1>
        <p>Handy tips to help you decide how to move</p>
        <p>Your weight allowance is {props.entitlement} lbs.</p>
        <h2 className="sm-heading">One move, several parts</h2>
        <p>
          It’s common to move some things yourself and have professional movers do the rest. You can also move things to
          or from more than one location.
        </p>
        <h2 className="sm-heading">Keep important things with you</h2>
        <p>
          It’s a smart idea to move essential items, heirlooms, and irreplaceable things yourself. Select a PPM
          (personally procured move) to do that.
        </p>
        <h2 className="sm-heading">Spread out your pickup dates</h2>
        <p>
          The easiest way to get the right boxes to the right location is to request different pickup and delivery days
          for different loads.
        </p>
      </div>
    </div>
  </div>
);

MovingInfo.propTypes = {
  entitlement: PropTypes.string,
};

const mapStateToProps = (state) => {
  const entitlement = loadEntitlementsFromState(state);
  return {
    entitlement,
  };
};

export default connect(mapStateToProps)(MovingInfo);
