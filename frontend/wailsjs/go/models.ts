export namespace backend {
	
	export class MoveWinrate {
	    san: string;
	    uci: string;
	    total: number;
	    whiteRate: number;
	    blackRate: number;
	    drawRate: number;
	    chance: number;
	
	    static createFrom(source: any = {}) {
	        return new MoveWinrate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.san = source["san"];
	        this.uci = source["uci"];
	        this.total = source["total"];
	        this.whiteRate = source["whiteRate"];
	        this.blackRate = source["blackRate"];
	        this.drawRate = source["drawRate"];
	        this.chance = source["chance"];
	    }
	}
	export class PositionWinrate {
	    total: number;
	    whiteRate: number;
	    blackRate: number;
	    drawRate: number;
	    moves: MoveWinrate[];
	
	    static createFrom(source: any = {}) {
	        return new PositionWinrate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.whiteRate = source["whiteRate"];
	        this.blackRate = source["blackRate"];
	        this.drawRate = source["drawRate"];
	        this.moves = this.convertValues(source["moves"], MoveWinrate);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Repertoire {
	    id: number;
	    name: string;
	    color: string;
	    elo: number;
	    coverage: number;
	
	    static createFrom(source: any = {}) {
	        return new Repertoire(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	        this.elo = source["elo"];
	        this.coverage = source["coverage"];
	    }
	}

}

