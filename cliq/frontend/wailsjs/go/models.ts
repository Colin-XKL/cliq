export namespace main {
	
	export class Variable {
	    type: string;
	    arg_name?: string;
	    label: string;
	    description: string;
	    required: boolean;
	    options?: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Variable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.arg_name = source["arg_name"];
	        this.label = source["label"];
	        this.description = source["description"];
	        this.required = source["required"];
	        this.options = source["options"];
	    }
	}
	export class Command {
	    id: string;
	    name: string;
	    description: string;
	    command: string;
	    variables: Record<string, Variable>;
	
	    static createFrom(source: any = {}) {
	        return new Command(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.command = source["command"];
	        this.variables = this.convertValues(source["variables"], Variable, true);
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
	export class TemplateFile {
	    name: string;
	    description: string;
	    version: string;
	    author: string;
	    cliq_template_version: string;
	    cmds: Command[];
	
	    static createFrom(source: any = {}) {
	        return new TemplateFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.version = source["version"];
	        this.author = source["author"];
	        this.cliq_template_version = source["cliq_template_version"];
	        this.cmds = this.convertValues(source["cmds"], Command);
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

}

