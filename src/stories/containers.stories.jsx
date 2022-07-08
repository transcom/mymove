import React from 'react';

// Containers
export default {
  title: 'Components/Containers',
};

export const all = () => (
  <div id="containers" style={{ padding: '20px' }}>
    <div className="container">
      <code>
        <b>Container Default</b>
        <br />
        .container
      </code>
    </div>
    <div className="container container--gray">
      <code>
        <b>Container Gray</b>
        <br />
        .container
        <br />
        .container--gray
      </code>
    </div>
    <div className="container container--popout">
      <code>
        <b>Container Popout</b>
        <br />
        .container
        <br />
        .container--popout
      </code>
    </div>
    <div className="container container--accent--hhg">
      <code>
        <b>Container Accent HHG</b>
        <br />
        .container
        <br />
        .container--accent--hhg
      </code>
    </div>
    <div className="container container--accent--ppm">
      <code>
        <b>Container Accent PPM</b>
        <br />
        .container
        <br />
        .container--accent--ppm
      </code>
    </div>
    <div className="container container--accent--ub">
      <code>
        <b>Container Accent UB</b>
        <br />
        .container
        <br />
        .container--accent--ub
      </code>
    </div>
    <div className="container container--accent--nts">
      <code>
        <b>Container Accent NTS</b>
        <br />
        .container
        <br />
        .container--accent--nts
      </code>
    </div>
    <div className="container container--accent--ntsr">
      <code>
        <b>Container Accent NTSR</b>
        <br />
        .container
        <br />
        .container--accent--ntsr
      </code>
    </div>
  </div>
);
