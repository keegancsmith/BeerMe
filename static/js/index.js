var BeerMe = React.createClass({
  handleDrinkSelect: function(drink) {
    state = this.state;
    console.assert(state.drink === null && state.team === null);
    state.drink = drink;
    this.setState(state);
  },
  handleTeamSelect: function(drink) {
    state = this.state;
    console.assert(state.drink !== null && state.team === null);
    state.team = drink;
    state.submitting = true;
    this.setState(state);
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      type: 'POST',
      data: JSON.stringify(state),
      success: function() {
        state.submitting = false;
        this.setState(state)
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {drink: null, team: null, submitting: false};
  },
  render: function() {
    if (this.state.drink === null) {
      return (
        <SimpleButtonGroup
          labels={this.props.drinks}
          onClick={this.handleDrinkSelect}
        />
      );
    } else if (this.state.team === null) {
      return (
        <SimpleButtonGroup
          labels={this.props.teams}
          onClick={this.handleTeamSelect}
        />
      );
    } else if (this.state.submitting) {
      return <div className="spinner-loader">Submitting...</div>;
    } else {
      return <h1>{this.state.drink} - {this.state.team}</h1>;
    }
  }
});

var SimpleButtonGroup = React.createClass({
  render: function() {
    var onClick = this.props.onClick;
    var nodes = this.props.labels.map(function (label) {
      var handler = function() { onClick(label) };
      return (
        <button type="button" className="btn btn-default" onClick={handler}>
          {label}
        </button>
      );
    });
    return (
      <div className="btn-group-vertical" role="group">
        {nodes}
      </div>
    );
  }
});

var drinks = ["Amstel", "Black Label", "Castle"]
var teams = ["Business Optics", "Siyelo", "Sourcegraph"]

React.render(
  <div className="well center-block">
    <BeerMe drinks={drinks} teams={teams} url="suip" />
  </div>,
  document.getElementById('content')
);
